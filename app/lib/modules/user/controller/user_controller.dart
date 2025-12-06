import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../model/user_info.dart';
import '../service/user_service.dart';

final userServiceProvider = Provider((ref) => UserService());

final userInfoProvider = FutureProvider<UserInfo>((ref) async {
  final userService = ref.read(userServiceProvider);
  return await userService.getUserInfo();
});
