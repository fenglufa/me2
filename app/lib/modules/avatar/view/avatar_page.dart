import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:go_router/go_router.dart';
import 'package:cached_network_image/cached_network_image.dart';
import 'package:image_picker/image_picker.dart';
import '../controller/avatar_controller.dart';
import '../../../core/utils/image_picker_service.dart';
import '../service/avatar_oss_upload_service.dart';
import '../../diary/controller/diary_controller.dart';

class AvatarPage extends ConsumerStatefulWidget {
  const AvatarPage({super.key});

  @override
  ConsumerState<AvatarPage> createState() => _AvatarPageState();
}

class _AvatarPageState extends ConsumerState<AvatarPage> {
  late TextEditingController _nameController;
  late TextEditingController _occupationController;
  XFile? _selectedImage;
  bool _isCreating = false;

  // Form fields
  int _gender = 1; // 1:男 2:女 3:其他
  DateTime? _birthDate;
  int _maritalStatus = 1; // 1:单身 2:恋爱中 3:已婚 4:其他

  @override
  void initState() {
    super.initState();
    _nameController = TextEditingController();
    _occupationController = TextEditingController();
  }

  @override
  void dispose() {
    _nameController.dispose();
    _occupationController.dispose();
    super.dispose();
  }

  Future<void> _pickImage() async {
    final imagePickerService = ref.read(imagePickerServiceProvider);
    final image = await imagePickerService.showImagePickerBottomSheet(context);
    if (image != null) {
      setState(() {
        _selectedImage = image;
      });
    }
  }

  Future<void> _createAvatar() async {
    if (_nameController.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入分身名称')),
      );
      return;
    }

    if (_birthDate == null) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请选择出生日期')),
      );
      return;
    }

    if (_occupationController.text.trim().isEmpty) {
      ScaffoldMessenger.of(context).showSnackBar(
        const SnackBar(content: Text('请输入职业')),
      );
      return;
    }

    if (_isCreating) return;

    setState(() {
      _isCreating = true;
    });

    try {
      final avatarService = ref.read(avatarServiceProvider);

      // Format birth date as YYYY-MM-DD
      final birthDateStr = '${_birthDate!.year.toString().padLeft(4, '0')}-'
          '${_birthDate!.month.toString().padLeft(2, '0')}-'
          '${_birthDate!.day.toString().padLeft(2, '0')}';

      // Step 1: Create avatar without avatar URL first
      final createdAvatar = await avatarService.createAvatar(
        name: _nameController.text.trim(),
        avatarUrl: null,
        gender: _gender,
        birthDate: birthDateStr,
        occupation: _occupationController.text.trim(),
        maritalStatus: _maritalStatus,
      );

      // Step 2: Upload avatar if selected, then update
      if (_selectedImage != null) {
        final ossService = ref.read(avatarOssUploadServiceProvider);
        final avatarUrl = await ossService.uploadAvatarAvatar(File(_selectedImage!.path));

        // Update avatar with the uploaded avatar URL
        await avatarService.updateAvatar(
          id: createdAvatar.id,
          avatarUrl: avatarUrl,
        );
      }

      // Refresh avatar info
      ref.invalidate(myAvatarProvider);

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('分身创建成功')),
        );
        setState(() {
          _nameController.clear();
          _occupationController.clear();
          _selectedImage = null;
          _gender = 1;
          _birthDate = null;
          _maritalStatus = 1;
        });
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('创建失败: $e')),
        );
      }
    } finally {
      if (mounted) {
        setState(() {
          _isCreating = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final avatarAsync = ref.watch(myAvatarProvider);

    return Scaffold(
      body: SafeArea(
        child: avatarAsync.when(
          data: (avatar) => SingleChildScrollView(
            child: Column(
              children: [
                _buildHeader(context, avatar),
                const SizedBox(height: 16),
                _buildDiarySection(context),
                const SizedBox(height: 16),
                _buildGrowthSection(context),
              ],
            ),
          ),
          loading: () => const Center(child: CircularProgressIndicator()),
          error: (error, _) {
            // Check if error is "no avatar created yet"
            if (error.toString().contains('还没有创建分身')) {
              return _buildCreateAvatarView();
            }
            // Other errors
            return Center(
              child: Column(
                mainAxisAlignment: MainAxisAlignment.center,
                children: [
                  Text('加载失败: $error'),
                  const SizedBox(height: 16),
                  ElevatedButton(
                    onPressed: () => ref.invalidate(myAvatarProvider),
                    child: const Text('重试'),
                  ),
                ],
              ),
            );
          },
        ),
      ),
    );
  }

  Widget _buildCreateAvatarView() {
    return SingleChildScrollView(
      padding: const EdgeInsets.all(24),
      child: Column(
        mainAxisAlignment: MainAxisAlignment.center,
        children: [
          const SizedBox(height: 60),
          const Icon(
            Icons.person_add_outlined,
            size: 80,
            color: Colors.purple,
          ),
          const SizedBox(height: 24),
          const Text(
            '创建你的分身',
            style: TextStyle(
              fontSize: 28,
              fontWeight: FontWeight.bold,
            ),
          ),
          const SizedBox(height: 12),
          Text(
            '为你的 AI 分身取一个名字，开始你的成长之旅',
            style: TextStyle(
              fontSize: 16,
              color: Colors.grey.shade600,
            ),
            textAlign: TextAlign.center,
          ),
          const SizedBox(height: 48),
          GestureDetector(
            onTap: _pickImage,
            child: Stack(
              children: [
                CircleAvatar(
                  radius: 60,
                  backgroundColor: Colors.purple.shade100,
                  backgroundImage: _selectedImage != null
                      ? FileImage(File(_selectedImage!.path))
                      : null,
                  child: _selectedImage == null
                      ? const Icon(Icons.person, size: 60, color: Colors.purple)
                      : null,
                ),
                Positioned(
                  right: 0,
                  bottom: 0,
                  child: Container(
                    padding: const EdgeInsets.all(8),
                    decoration: BoxDecoration(
                      color: Colors.purple,
                      shape: BoxShape.circle,
                      border: Border.all(color: Colors.white, width: 2),
                    ),
                    child: const Icon(
                      Icons.camera_alt,
                      size: 20,
                      color: Colors.white,
                    ),
                  ),
                ),
              ],
            ),
          ),
          const SizedBox(height: 12),
          Text(
            '点击上传分身头像（可选）',
            style: TextStyle(
              fontSize: 14,
              color: Colors.grey.shade600,
            ),
          ),
          const SizedBox(height: 32),
          TextField(
            controller: _nameController,
            decoration: InputDecoration(
              labelText: '分身名称*',
              hintText: '给你的分身起个名字',
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              prefixIcon: const Icon(Icons.badge_outlined),
            ),
            textInputAction: TextInputAction.next,
          ),
          const SizedBox(height: 16),
          DropdownButtonFormField<int>(
            value: _gender,
            decoration: InputDecoration(
              labelText: '性别*',
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              prefixIcon: const Icon(Icons.person_outline),
            ),
            items: const [
              DropdownMenuItem(value: 1, child: Text('男')),
              DropdownMenuItem(value: 2, child: Text('女')),
              DropdownMenuItem(value: 3, child: Text('其他')),
            ],
            onChanged: (value) {
              if (value != null) {
                setState(() {
                  _gender = value;
                });
              }
            },
          ),
          const SizedBox(height: 16),
          InkWell(
            onTap: () async {
              final picked = await showDatePicker(
                context: context,
                initialDate: _birthDate ?? DateTime(2000, 1, 1),
                firstDate: DateTime(1950),
                lastDate: DateTime.now(),
              );
              if (picked != null) {
                setState(() {
                  _birthDate = picked;
                });
              }
            },
            child: InputDecorator(
              decoration: InputDecoration(
                labelText: '出生日期*',
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
                prefixIcon: const Icon(Icons.calendar_today_outlined),
              ),
              child: Text(
                _birthDate == null
                    ? '请选择出生日期'
                    : '${_birthDate!.year}年${_birthDate!.month}月${_birthDate!.day}日',
                style: TextStyle(
                  color: _birthDate == null ? Colors.grey : Colors.black,
                ),
              ),
            ),
          ),
          const SizedBox(height: 16),
          TextField(
            controller: _occupationController,
            decoration: InputDecoration(
              labelText: '职业*',
              hintText: '例如：学生、工程师、医生等',
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              prefixIcon: const Icon(Icons.work_outline),
            ),
            textInputAction: TextInputAction.next,
          ),
          const SizedBox(height: 16),
          DropdownButtonFormField<int>(
            value: _maritalStatus,
            decoration: InputDecoration(
              labelText: '婚姻状态*',
              border: OutlineInputBorder(
                borderRadius: BorderRadius.circular(12),
              ),
              prefixIcon: const Icon(Icons.favorite_outline),
            ),
            items: const [
              DropdownMenuItem(value: 1, child: Text('单身')),
              DropdownMenuItem(value: 2, child: Text('恋爱中')),
              DropdownMenuItem(value: 3, child: Text('已婚')),
              DropdownMenuItem(value: 4, child: Text('其他')),
            ],
            onChanged: (value) {
              if (value != null) {
                setState(() {
                  _maritalStatus = value;
                });
              }
            },
          ),
          const SizedBox(height: 32),
          SizedBox(
            width: double.infinity,
            height: 50,
            child: ElevatedButton(
              onPressed: _isCreating ? null : _createAvatar,
              style: ElevatedButton.styleFrom(
                backgroundColor: Colors.purple,
                foregroundColor: Colors.white,
                shape: RoundedRectangleBorder(
                  borderRadius: BorderRadius.circular(12),
                ),
              ),
              child: _isCreating
                  ? const SizedBox(
                      width: 24,
                      height: 24,
                      child: CircularProgressIndicator(
                        strokeWidth: 2,
                        valueColor: AlwaysStoppedAnimation(Colors.white),
                      ),
                    )
                  : const Text(
                      '创建分身',
                      style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
                    ),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildHeader(BuildContext context, dynamic avatar) {
    return Container(
      padding: const EdgeInsets.all(20),
      child: Row(
        children: [
          CircleAvatar(
            radius: 40,
            backgroundColor: Colors.purple.shade100,
            backgroundImage: avatar.avatarUrl.isNotEmpty
                ? CachedNetworkImageProvider(avatar.avatarUrl)
                : null,
            child: avatar.avatarUrl.isEmpty
                ? const Icon(Icons.person, size: 40)
                : null,
          ),
          const SizedBox(width: 16),
          Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  avatar.name,
                  style: const TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
                ),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(Icons.edit_outlined),
            onPressed: () {
              context.push('/edit-avatar', extra: avatar.id);
            },
          ),
        ],
      ),
    );
  }

  Widget _buildDiarySection(BuildContext context) {
    final diariesAsync = ref.watch(avatarDiariesProvider);

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                '分身日记',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              TextButton(
                onPressed: () {
                  context.push('/avatar-diaries');
                },
                child: const Text('查看全部'),
              ),
            ],
          ),
          diariesAsync.when(
            data: (response) {
              if (response.list.isEmpty) {
                return const Padding(
                  padding: EdgeInsets.all(16),
                  child: Text('暂无日记', style: TextStyle(color: Colors.grey)),
                );
              }
              final diaries = response.list.take(2).toList();
              return Column(
                children: diaries.map((diary) => _buildDiaryCard(
                  diary.title,
                  diary.content,
                  diary.date,
                )).toList(),
              );
            },
            loading: () => const Center(child: CircularProgressIndicator()),
            error: (error, _) => Padding(
              padding: const EdgeInsets.all(16),
              child: Text('加载失败: $error', style: const TextStyle(color: Colors.red)),
            ),
          ),
        ],
      ),
    );
  }

  Widget _buildDiaryCard(String title, String content, String time) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 8),
            Text(content, style: const TextStyle(color: Colors.grey)),
            const SizedBox(height: 8),
            Text(time, style: TextStyle(fontSize: 12, color: Colors.grey.shade600)),
          ],
        ),
      ),
    );
  }

  Widget _buildGrowthSection(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '成长节点',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 12),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: [
              _buildMilestoneChip('第一次旅行', true),
              _buildMilestoneChip('第一次社交', true),
              _buildMilestoneChip('第一次创作', false),
              _buildMilestoneChip('第一次探索', true),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildMilestoneChip(String label, bool achieved) {
    return Chip(
      label: Text(label),
      avatar: Icon(
        achieved ? Icons.check_circle : Icons.circle_outlined,
        size: 16,
        color: achieved ? Colors.green : Colors.grey,
      ),
      backgroundColor: achieved ? Colors.green.shade50 : Colors.grey.shade100,
    );
  }
}
