class Event {
  final int id;
  final int avatarId;
  final String eventType;
  final String eventTitle;
  final String eventText;
  final String imageUrl;
  final int sceneId;
  final String sceneName;
  final int occurredAt;

  Event({
    required this.id,
    required this.avatarId,
    required this.eventType,
    required this.eventTitle,
    required this.eventText,
    required this.imageUrl,
    required this.sceneId,
    required this.sceneName,
    required this.occurredAt,
  });

  factory Event.fromJson(Map<String, dynamic> json) => Event(
        id: json['id'],
        avatarId: json['avatar_id'],
        eventType: json['event_type'],
        eventTitle: json['event_title'] ?? '',
        eventText: json['event_text'],
        imageUrl: json['image_url'] ?? '',
        sceneId: json['scene_id'],
        sceneName: json['scene_name'],
        occurredAt: json['occurred_at'],
      );
}
