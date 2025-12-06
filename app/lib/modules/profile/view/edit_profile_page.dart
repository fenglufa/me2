import 'dart:io';
import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:image_picker/image_picker.dart';
import '../../../core/utils/image_picker_service.dart';
import '../../../core/utils/oss_upload_service.dart';
import '../../user/controller/user_controller.dart';
import '../../user/service/user_service.dart';

class EditProfilePage extends ConsumerStatefulWidget {
  const EditProfilePage({super.key});

  @override
  ConsumerState<EditProfilePage> createState() => _EditProfilePageState();
}

class _EditProfilePageState extends ConsumerState<EditProfilePage> {
  late TextEditingController _nicknameController;
  XFile? _selectedImage;
  bool _isLoading = false;

  @override
  void initState() {
    super.initState();
    _nicknameController = TextEditingController();

    WidgetsBinding.instance.addPostFrameCallback((_) {
      final userInfo = ref.read(userInfoProvider).value;
      if (userInfo != null) {
        _nicknameController.text = userInfo.nickname;
      }
    });
  }

  @override
  void dispose() {
    _nicknameController.dispose();
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
        final ossService = ref.read(ossUploadServiceProvider);
        avatarUrl = await ossService.uploadAvatar(File(_selectedImage!.path));
      }

      // Update user info
      final userService = UserService();
      await userService.updateUser(
        nickname: _nicknameController.text.trim().isEmpty ? null : _nicknameController.text.trim(),
        avatarUrl: avatarUrl,
      );

      // Refresh user info
      ref.invalidate(userInfoProvider);

      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          const SnackBar(content: Text('保存成功')),
        );
        Navigator.pop(context);
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
    final userInfoAsync = ref.watch(userInfoProvider);

    return Scaffold(
      appBar: AppBar(
        title: const Text('编辑资料'),
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
      body: userInfoAsync.when(
        data: (userInfo) => SingleChildScrollView(
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
                          : (userInfo.avatarUrl.isNotEmpty
                              ? NetworkImage(userInfo.avatarUrl)
                              : null) as ImageProvider?,
                      child: _selectedImage == null && userInfo.avatarUrl.isEmpty
                          ? const Icon(Icons.person, size: 50)
                          : null,
                    ),
                    Positioned(
                      right: 0,
                      bottom: 0,
                      child: Container(
                        padding: const EdgeInsets.all(4),
                        decoration: BoxDecoration(
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
                controller: _nicknameController,
                decoration: const InputDecoration(
                  labelText: '昵称',
                  border: OutlineInputBorder(),
                ),
              ),
            ],
          ),
        ),
        loading: () => const Center(child: CircularProgressIndicator()),
        error: (error, _) => Center(child: Text('加载失败: $error')),
      ),
    );
  }
}
