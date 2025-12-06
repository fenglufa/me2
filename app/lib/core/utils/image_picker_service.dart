import 'package:flutter/material.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:image_picker/image_picker.dart';
import '../permissions/permission_service.dart';

class ImagePickerService {
  final ImagePicker _picker = ImagePicker();
  final PermissionService _permissionService;

  ImagePickerService(this._permissionService);

  Future<XFile?> pickFromGallery(BuildContext context) async {
    bool hasPermission = await _permissionService.getPhotosPermission(context);
    if (!hasPermission) return null;

    try {
      return await _picker.pickImage(source: ImageSource.gallery, imageQuality: 80);
    } catch (e) {
      return null;
    }
  }

  Future<XFile?> takePhoto(BuildContext context) async {
    bool hasPermission = await _permissionService.getCameraPermission(context);
    if (!hasPermission) return null;

    try {
      return await _picker.pickImage(source: ImageSource.camera, imageQuality: 80);
    } catch (e) {
      return null;
    }
  }

  Future<XFile?> showImagePickerBottomSheet(BuildContext context) async {
    final String? choice = await showModalBottomSheet<String>(
      context: context,
      shape: const RoundedRectangleBorder(
        borderRadius: BorderRadius.vertical(top: Radius.circular(20)),
      ),
      builder: (context) => SafeArea(
        child: Column(
          mainAxisSize: MainAxisSize.min,
          children: [
            Container(
              width: 40,
              height: 4,
              margin: const EdgeInsets.symmetric(vertical: 12),
              decoration: BoxDecoration(
                color: Colors.grey.shade300,
                borderRadius: BorderRadius.circular(2),
              ),
            ),
            Padding(
              padding: const EdgeInsets.all(16),
              child: Column(
                children: [
                  ListTile(
                    leading: const Icon(Icons.camera_alt),
                    title: const Text('拍照'),
                    onTap: () => Navigator.pop(context, 'camera'),
                  ),
                  ListTile(
                    leading: const Icon(Icons.photo_library),
                    title: const Text('从相册选择'),
                    onTap: () => Navigator.pop(context, 'gallery'),
                  ),
                  const SizedBox(height: 8),
                  SizedBox(
                    width: double.infinity,
                    child: TextButton(
                      onPressed: () => Navigator.pop(context),
                      child: const Text('取消', style: TextStyle(color: Colors.grey)),
                    ),
                  ),
                ],
              ),
            ),
          ],
        ),
      ),
    );

    if (context.mounted && choice == 'camera') {
      return await takePhoto(context);
    } else if (context.mounted && choice == 'gallery') {
      return await pickFromGallery(context);
    }

    return null;
  }
}

final imagePickerServiceProvider = Provider<ImagePickerService>((ref) {
  final permissionService = ref.read(permissionServiceProvider);
  return ImagePickerService(permissionService);
});
