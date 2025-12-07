import 'package:dio/dio.dart';
import '../../../core/network/dio_client.dart';
import '../../../core/constants/api_endpoints.dart';
import '../model/avatar_info.dart';

class AvatarService {
  final Dio _dio = DioClient.instance;

  Future<AvatarInfo> getMyAvatar() async {
    final response = await _dio.get(ApiEndpoints.myAvatar);
    if (response.data == null) {
      throw Exception('您还没有创建分身，请先创建分身');
    }
    return AvatarInfo.fromJson(response.data);
  }

  Future<AvatarInfo> getAvatar(int id) async {
    final response = await _dio.get('${ApiEndpoints.avatarDetail}/$id');
    return AvatarInfo.fromJson(response.data);
  }

  Future<AvatarInfo> createAvatar({
    required String name,
    String? avatarUrl,
    required int gender,
    required String birthDate,
    required String occupation,
    required int maritalStatus,
  }) async {
    final response = await _dio.post(
      ApiEndpoints.createAvatar,
      data: {
        'name': name,
        'avatar_url': avatarUrl ?? '',
        'gender': gender,
        'birth_date': birthDate,
        'occupation': occupation,
        'marital_status': maritalStatus,
      },
    );
    return AvatarInfo.fromJson(response.data);
  }

  Future<AvatarInfo> updateAvatar({
    required int id,
    String? name,
    String? avatarUrl,
  }) async {
    final response = await _dio.patch(
      '${ApiEndpoints.avatarDetail}/$id',
      data: {
        if (name != null) 'name': name,
        if (avatarUrl != null) 'avatar_url': avatarUrl,
      },
    );
    return AvatarInfo.fromJson(response.data);
  }
}
