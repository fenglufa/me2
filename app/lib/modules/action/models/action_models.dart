class ActionLog {
  final int id;
  final int avatarId;
  final String actionType;
  final int sceneId;
  final String sceneName;
  final double intentScore;
  final String triggerReason;
  final int createdAt;

  ActionLog({
    required this.id,
    required this.avatarId,
    required this.actionType,
    required this.sceneId,
    required this.sceneName,
    required this.intentScore,
    required this.triggerReason,
    required this.createdAt,
  });

  factory ActionLog.fromJson(Map<String, dynamic> json) => ActionLog(
        id: json['id'],
        avatarId: json['avatar_id'],
        actionType: json['action_type'],
        sceneId: json['scene_id'],
        sceneName: json['scene_name'],
        intentScore: (json['intent_score'] as num).toDouble(),
        triggerReason: json['trigger_reason'] ?? '',
        createdAt: json['created_at'],
      );
}
