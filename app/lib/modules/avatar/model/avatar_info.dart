class AvatarInfo {
  final int id;
  final int userId;
  final String name;
  final String avatarUrl;
  final int gender;
  final String birthDate;
  final String occupation;
  final int maritalStatus;
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
      'created_at': createdAt,
    };
  }
}
