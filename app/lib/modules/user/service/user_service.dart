import 'dart:io';
import 'package:dio/dio.dart';
import '../../../core/network/dio_client.dart';
import '../../../core/constants/api_endpoints.dart';
import '../model/user_info.dart';

class UserService {
  final Dio _dio = DioClient.instance;

  Future<UserInfo> getUserInfo() async {
    final response = await _dio.get(ApiEndpoints.userInfo);
    return UserInfo.fromJson(response.data);
  }

  Future<void> updateUser({String? nickname, String? avatarUrl}) async {
    await _dio.put(
      ApiEndpoints.updateUser,
      data: {
        if (nickname != null) 'nickname': nickname,
        if (avatarUrl != null) 'avatar_url': avatarUrl,
      },
    );
  }

  Future<String> uploadAvatar(File file) async {
    // Get upload token
    final tokenResponse = await _dio.get(ApiEndpoints.avatarUploadToken);
    final uploadUrl = tokenResponse.data['upload_url'];
    final key = tokenResponse.data['key'];

    // Upload to OSS
    final formData = FormData.fromMap({
      'file': await MultipartFile.fromFile(file.path),
    });
    await Dio().put(uploadUrl, data: formData);

    // Complete upload
    final completeResponse = await _dio.post(
      ApiEndpoints.completeAvatarUpload,
      data: {'key': key},
    );
    return completeResponse.data['avatar_url'];
  }
}
