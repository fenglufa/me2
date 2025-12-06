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
}
