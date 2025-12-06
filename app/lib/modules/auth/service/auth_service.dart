import 'package:dio/dio.dart';
import '../../../core/network/dio_client.dart';
import '../../../core/constants/api_endpoints.dart';
import '../../../core/storage/token_storage.dart';
import '../model/login_request.dart';
import '../model/login_response.dart';

class AuthService {
  final Dio _dio = DioClient.instance;

  Future<void> sendCode(String phone) async {
    await _dio.post(
      ApiEndpoints.sendCode,
      data: {'phone': phone},
    );
  }

  Future<LoginResponse> login(String phone, String code) async {
    final response = await _dio.post(
      ApiEndpoints.login,
      data: LoginRequest(phone: phone, code: code).toJson(),
    );

    final loginResponse = LoginResponse.fromJson(response.data);
    await TokenStorage.saveToken(loginResponse.token);
    await TokenStorage.saveUserId(loginResponse.userId);

    return loginResponse;
  }

  Future<void> logout() async {
    await TokenStorage.clear();
  }

  Future<bool> isLoggedIn() async {
    final token = await TokenStorage.getToken();
    return token != null;
  }
}
