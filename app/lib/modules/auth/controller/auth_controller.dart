import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../service/auth_service.dart';
import '../model/login_response.dart';

final authServiceProvider = Provider((ref) => AuthService());

final authControllerProvider =
    StateNotifierProvider<AuthController, AsyncValue<LoginResponse?>>((ref) {
  return AuthController(ref.read(authServiceProvider));
});

class AuthController extends StateNotifier<AsyncValue<LoginResponse?>> {
  final AuthService _authService;

  AuthController(this._authService) : super(const AsyncValue.data(null));

  Future<void> sendCode(String phone) async {
    try {
      await _authService.sendCode(phone);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
    }
  }

  Future<void> login(String phone, String code) async {
    state = const AsyncValue.loading();
    try {
      final response = await _authService.login(phone, code);
      state = AsyncValue.data(response);
    } catch (e, st) {
      state = AsyncValue.error(e, st);
    }
  }

  Future<void> logout() async {
    await _authService.logout();
    state = const AsyncValue.data(null);
  }

  Future<bool> checkLoginStatus() async {
    return await _authService.isLoggedIn();
  }
}
