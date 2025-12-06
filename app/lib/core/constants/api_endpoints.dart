class ApiEndpoints {
  // Auth
  static const String sendCode = '/api/v1/user/send-code';
  static const String login = '/api/v1/user/login';

  // User
  static const String userInfo = '/api/v1/user/info';
  static const String updateUser = '/api/v1/user/update';
  static const String subscription = '/api/v1/user/subscription';

  // Avatar
  static const String createAvatar = '/api/v1/avatar/create';
  static const String myAvatar = '/api/v1/avatar/my';
  static const String avatarDetail = '/api/v1/avatar';
  static const String avatarUploadToken = '/api/v1/user/avatar-token';
  static const String completeAvatarUpload = '/api/v1/user/avatar-complete';

  // World
  static const String maps = '/api/v1/world/maps';
  static const String regions = '/api/v1/world/regions';
  static const String scenes = '/api/v1/world/scenes';

  // Event
  static const String timeline = '/api/v1/events/timeline';

  // Action
  static const String lastAction = '/api/v1/action/last';
  static const String actionHistory = '/api/v1/action/history';

  // Diary
  static const String avatarDiaries = '/api/v1/diary/avatar';
  static const String userDiaries = '/api/v1/diary/user';
  static const String diaryStats = '/api/v1/diary/stats';
}
