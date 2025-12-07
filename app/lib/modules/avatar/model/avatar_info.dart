class AvatarInfo {
  final int id;
  final int userId;
  final String name;
  final String avatarUrl;
  final int gender;
  final String birthDate;
  final String occupation;
  final int maritalStatus;
  final double warmth;
  final double adventurous;
  final double social;
  final double creative;
  final double calm;
  final double energetic;
  final int createdAt;

  AvatarInfo({
    required this.id,
    required this.userId,
    required this.name,
    required this.avatarUrl,
    required this.gender,
    required this.birthDate,
    required this.occupation,
    required this.maritalStatus,
    required this.warmth,
    required this.adventurous,
    required this.social,
    required this.creative,
    required this.calm,
    required this.energetic,
    required this.createdAt,
  });

  factory AvatarInfo.fromJson(Map<String, dynamic> json) {
    return AvatarInfo(
      id: json['id'] ?? 0,
      userId: json['user_id'] ?? 0,
      name: json['name'] ?? '',
      avatarUrl: json['avatar_url'] ?? '',
      gender: json['gender'] ?? 1,
      birthDate: json['birth_date'] ?? '',
      occupation: json['occupation'] ?? '',
      maritalStatus: json['marital_status'] ?? 1,
      warmth: (json['warmth'] ?? 0.0).toDouble(),
      adventurous: (json['adventurous'] ?? 0.0).toDouble(),
      social: (json['social'] ?? 0.0).toDouble(),
      creative: (json['creative'] ?? 0.0).toDouble(),
      calm: (json['calm'] ?? 0.0).toDouble(),
      energetic: (json['energetic'] ?? 0.0).toDouble(),
      createdAt: json['created_at'] ?? 0,
    );
  }

  Map<String, dynamic> toJson() {
    return {
      'id': id,
      'user_id': userId,
      'name': name,
      'avatar_url': avatarUrl,
      'gender': gender,
      'birth_date': birthDate,
      'occupation': occupation,
      'marital_status': maritalStatus,
      'warmth': warmth,
      'adventurous': adventurous,
      'social': social,
      'creative': creative,
      'calm': calm,
      'energetic': energetic,
      'created_at': createdAt,
    };
  }
}
