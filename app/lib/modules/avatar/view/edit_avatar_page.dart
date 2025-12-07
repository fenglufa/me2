import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:image_picker/image_picker.dart';
import 'package:go_router/go_router.dart';
import '../../../core/utils/image_picker_service.dart';
import '../service/avatar_oss_upload_service.dart';
import '../controller/avatar_controller.dart';

class EditAvatarPage extends ConsumerStatefulWidget {
  final int avatarId;

  const EditAvatarPage({super.key, required this.avatarId});

  @override
  ConsumerState<EditAvatarPage> createState() => _EditAvatarPageState();
}

class _EditAvatarPageState extends ConsumerState<EditAvatarPage> {
  late TextEditingController _nameController;
  late TextEditingController _occupationController;
  XFile? _selectedImage;
  bool _isLoading = false;
  int? _gender;
  DateTime? _birthDate;
  int? _maritalStatus;

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

  Future<void> _save() async {
    if (_isLoading) return;

    setState(() {
      _isLoading = true;
    });

    try {
      String? avatarUrl;

      // Upload avatar if selected
      if (_selectedImage != null) {
        final ossService = ref.read(avatarOssUploadServiceProvider);
        avatarUrl = await ossService.uploadAvatarAvatar(File(_selectedImage!.path));
      }

      // Update avatar info
      final avatarService = ref.read(avatarServiceProvider);
      final birthDateStr = _birthDate != null
          ? '${_birthDate!.year.toString().padLeft(4, '0')}-'
              '${_birthDate!.month.toString().padLeft(2, '0')}-'
              '${_birthDate!.day.toString().padLeft(2, '0')}'
          : null;

      await avatarService.updateAvatar(
        id: widget.avatarId,
        name: _nameController.text.trim().isEmpty ? null : _nameController.text.trim(),
        avatarUrl: avatarUrl,
        gender: _gender,
        birthDate: birthDateStr,
        occupation: _occupationController.text.trim().isEmpty ? null : _occupationController.text.trim(),
        maritalStatus: _maritalStatus,
      );

      // Refresh avatar info
      ref.invalidate(myAvatarProvider);
      ref.invalidate(avatarProvider(widget.avatarId));

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('保存成功')),
        );
        context.pop();
      }
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('保存失败: $e')),
        );
      }
    } finally {
      if (mounted) {
        setState(() {
          _isLoading = false;
        });
      }
    }
  }

  @override
  Widget build(BuildContext context) {
    final avatarAsync = ref.watch(avatarProvider(widget.avatarId));

    return Scaffold(
      appBar: AppBar(
        title: const Text('编辑分身资料'),
        actions: [
          TextButton(
            onPressed: _isLoading ? null : _save,
            child: _isLoading
                ? const SizedBox(
                    width: 20,
                    height: 20,
                    child: CircularProgressIndicator(strokeWidth: 2),
                  )
                : const Text('保存'),
          ),
        ],
      ),
      body: avatarAsync.when(
        data: (avatar) {
          if (_nameController.text.isEmpty) {
            _nameController.text = avatar.name;
            _occupationController.text = avatar.occupation;
            _gender = avatar.gender;
            _maritalStatus = avatar.maritalStatus;
            if (avatar.birthDate.isNotEmpty) {
              final parts = avatar.birthDate.split('-');
              if (parts.length == 3) {
                _birthDate = DateTime(int.parse(parts[0]), int.parse(parts[1]), int.parse(parts[2]));
              }
            }
          }

          return SingleChildScrollView(
            padding: const EdgeInsets.all(16),
            child: Column(
              children: [
                GestureDetector(
                  onTap: _pickImage,
                  child: Stack(
                    children: [
                      CircleAvatar(
                        radius: 50,
                        backgroundColor: Colors.grey.shade300,
                        backgroundImage: _selectedImage != null
                            ? FileImage(File(_selectedImage!.path))
                            : (avatar.avatarUrl.isNotEmpty
                                ? NetworkImage(avatar.avatarUrl)
                                : null) as ImageProvider?,
                        child: _selectedImage == null && avatar.avatarUrl.isEmpty
                            ? const Icon(Icons.person, size: 50)
                            : null,
                      ),
                      Positioned(
                        right: 0,
                        bottom: 0,
                        child: Container(
                          padding: const EdgeInsets.all(4),
                          decoration: const BoxDecoration(
                            color: Colors.blue,
                            shape: BoxShape.circle,
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
                const SizedBox(height: 32),
                TextField(
                  controller: _nameController,
                  decoration: const InputDecoration(
                    labelText: '分身名称',
                    border: OutlineInputBorder(),
                  ),
                ),
                const SizedBox(height: 16),
                DropdownButtonFormField<int>(
                  value: _gender,
                  decoration: const InputDecoration(
                    labelText: '性别',
                    border: OutlineInputBorder(),
                  ),
                  items: const [
                    DropdownMenuItem(value: 1, child: Text('男')),
                    DropdownMenuItem(value: 2, child: Text('女')),
                    DropdownMenuItem(value: 3, child: Text('其他')),
                  ],
                  onChanged: (value) => setState(() => _gender = value),
                ),
                const SizedBox(height: 16),
                InkWell(
                  onTap: () async {
                    final picked = await showDatePicker(
                      context: context,
                      initialDate: _birthDate ?? DateTime(2000),
                      firstDate: DateTime(1900),
                      lastDate: DateTime.now(),
                    );
                    if (picked != null) setState(() => _birthDate = picked);
                  },
                  child: InputDecorator(
                    decoration: const InputDecoration(
                      labelText: '出生日期',
                      border: OutlineInputBorder(),
                    ),
                    child: Text(_birthDate != null
                        ? '${_birthDate!.year}-${_birthDate!.month.toString().padLeft(2, '0')}-${_birthDate!.day.toString().padLeft(2, '0')}'
                        : '请选择出生日期'),
                  ),
                ),
                const SizedBox(height: 16),
                TextField(
                  controller: _occupationController,
                  decoration: const InputDecoration(
                    labelText: '职业',
                    border: OutlineInputBorder(),
                  ),
                ),
                const SizedBox(height: 16),
                DropdownButtonFormField<int>(
                  value: _maritalStatus,
                  decoration: const InputDecoration(
                    labelText: '婚姻状态',
                    border: OutlineInputBorder(),
                  ),
                  items: const [
                    DropdownMenuItem(value: 1, child: Text('单身')),
                    DropdownMenuItem(value: 2, child: Text('恋爱中')),
                    DropdownMenuItem(value: 3, child: Text('已婚')),
                    DropdownMenuItem(value: 4, child: Text('其他')),
                  ],
                  onChanged: (value) => setState(() => _maritalStatus = value),
                ),
              ],
            ),
          );
        },
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, _) => Center(child: Text('加载失败: $error')),
      ),
    );
  }
}
