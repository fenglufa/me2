import 'dart:io';
import 'package:dio/dio.dart';
import 'package:flutter_riverpod/flutter_riverpod.dart';
import 'package:http_parser/http_parser.dart';
import 'package:path/path.dart' as path;
import '../../../core/network/dio_client.dart';
import '../../../core/constants/api_endpoints.dart';

class AvatarOssUploadService {
  final Dio _dio = DioClient.instance;

  Future<String> uploadAvatarAvatar(File file) async {
    // 1. Get upload token
    final tokenResponse = await _dio.get(ApiEndpoints.avatarAvatarUploadToken);
    final data = tokenResponse.data;

    final host = data['upload_url'] as String;
    final accessid = data['accessid'] as String;
    final policy = data['policy'] as String;
    final signature = data['signature'] as String;
    final fileKey = data['key'] as String; // Use backend generated key
    final token = data['token'] as String;

    // 2. Upload to OSS
    final fileName = path.basename(file.path);
    final contentType = _getContentType(fileName);

    final formData = FormData.fromMap({
      'key': fileKey,
      'policy': policy,
      'OSSAccessKeyId': accessid,
      'signature': signature,
      'success_action_status': '200',
      'file': await MultipartFile.fromFile(
        file.path,
        filename: fileName,
        contentType: contentType,
      ),
    });

    final uploadDio = Dio();
    await uploadDio.post(
      host,
      data: formData,
      options: Options(
        headers: {'Content-Type': 'multipart/form-data'},
      ),
    );

    // 3. Complete upload with token
    final completeResponse = await _dio.post(
      ApiEndpoints.avatarCompleteAvatarUpload,
      data: {
        'key': fileKey,
        'token': token,
      },
    );

    return completeResponse.data['avatar_url'] as String;
  }

  MediaType _getContentType(String fileName) {
    final extension = path.extension(fileName).toLowerCase();
    switch (extension) {
      case '.jpg':
      case '.jpeg':
        return MediaType('image', 'jpeg');
      case '.png':
        return MediaType('image', 'png');
      case '.gif':
        return MediaType('image', 'gif');
      case '.webp':
        return MediaType('image', 'webp');
      default:
        return MediaType('image', 'jpeg');
    }
  }
}

final avatarOssUploadServiceProvider = Provider<AvatarOssUploadService>((ref) {
  return AvatarOssUploadService();
});
