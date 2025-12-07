import 'package:flutter_riverpod/flutter_riverpod.dart';
import '../model/avatar_info.dart';
import '../service/avatar_service.dart';

final avatarServiceProvider = Provider<AvatarService>((ref) {
  return AvatarService();
});

final myAvatarProvider = FutureProvider<AvatarInfo>((ref) async {
  final service = ref.read(avatarServiceProvider);
  return await service.getMyAvatar();
});

final avatarProvider = FutureProvider.family<AvatarInfo, int>((ref, id) async {
  final service = ref.read(avatarServiceProvider);
  return await service.getAvatar(id);
});
