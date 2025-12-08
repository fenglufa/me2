class Session {
  final int id;
  final int userId;
  final int avatarId;
  final String title;
  final String lastMessage;
  final int createdAt;
  final int updatedAt;

  Session({
    required this.id,
    required this.userId,
    required this.avatarId,
    required this.title,
    required this.lastMessage,
    required this.createdAt,
    required this.updatedAt,
  });

  factory Session.fromJson(Map<String, dynamic> json) {
    return Session(
      id: json['id'],
      userId: json['user_id'],
      avatarId: json['avatar_id'],
      title: json['title'],
      lastMessage: json['last_message'],
      createdAt: json['created_at'],
      updatedAt: json['updated_at'],
    );
  }
}

class Message {
  final int id;
  final int sessionId;
  final String role;
  final String content;
  final int createdAt;

  Message({
    required this.id,
    required this.sessionId,
    required this.role,
    required this.content,
    required this.createdAt,
  });

  factory Message.fromJson(Map<String, dynamic> json) {
    return Message(
      id: json['id'],
      sessionId: json['session_id'],
      role: json['role'],
      content: json['content'],
      createdAt: json['created_at'],
    );
  }
}
