import 'package:flutter/material.dart';
import 'package:permission_handler/permission_handler.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';

class PermissionService {
  Future<bool> getPhotosPermission(BuildContext context) async {
    try {
      PermissionStatus status = await Permission.photos.status;

      if (status.isGranted || status.isLimited) {
        return true;
      }

      if (status.isPermanentlyDenied) {
        if (context.mounted) {
          _showPermissionDeniedDialog(context, '相册权限', '需要相册权限来选择照片');
        }
        return false;
      }

      status = await Permission.photos.request();

      if (status.isGranted || status.isLimited) {
        return true;
      } else if (status.isPermanentlyDenied) {
        if (context.mounted) {
          _showPermissionDeniedDialog(context, '相册权限', '需要相册权限来选择照片');
        }
        return false;
      }

      return false;
    } catch (e) {
      return false;
    }
  }

  Future<bool> getCameraPermission(BuildContext context) async {
    try {
      PermissionStatus status = await Permission.camera.status;

      if (status.isGranted) {
        return true;
      }

      if (status.isPermanentlyDenied) {
        if (context.mounted) {
          _showPermissionDeniedDialog(context, '相机权限', '需要相机权限来拍摄照片');
        }
        return false;
      }

      status = await Permission.camera.request();

      if (status.isGranted) {
        return true;
      } else if (status.isPermanentlyDenied) {
        if (context.mounted) {
          _showPermissionDeniedDialog(context, '相机权限', '需要相机权限来拍摄照片');
        }
        return false;
      }

      return false;
    } catch (e) {
      return false;
    }
  }

  void _showPermissionDeniedDialog(BuildContext context, String permissionName, String message) {
    showDialog(
      context: context,
      builder: (context) => AlertDialog(
        title: Text('$permissionName未开启'),
        content: Text(message),
        actions: [
          TextButton(
            onPressed: () => Navigator.of(context).pop(),
            child: const Text('取消'),
          ),
          TextButton(
            onPressed: () {
              Navigator.of(context).pop();
              openAppSettings();
            },
            child: const Text('去设置'),
          ),
        ],
      ),
    );
  }
}

final permissionServiceProvider = Provider<PermissionService>((ref) {
  return PermissionService();
});
