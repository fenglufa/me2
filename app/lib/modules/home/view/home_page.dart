import 'package:flutter/material.dart';
import '../../event/models/event_models.dart';
import '../../event/services/event_service.dart';
import '../../action/models/action_models.dart';
import '../../action/services/action_service.dart';

class HomePage extends StatefulWidget {
  const HomePage({super.key});

  @override
  State<HomePage> createState() => _HomePageState();
}

class _HomePageState extends State<HomePage> {
  final _eventService = EventService();
  final _actionService = ActionService();
  List<Event> _recentEvents = [];
  ActionLog? _lastAction;
  bool _loadingEvents = false;
  bool _loadingAction = false;

  @override
  void initState() {
    super.initState();
    _loadRecentEvents();
    _loadLastAction();
  }

  Future<void> _loadRecentEvents() async {
    try {
      _loadingEvents = true;
      setState(() {});
      _recentEvents = await _eventService.getEventTimeline(pageSize: 10);
    } catch (e) {
      debugPrint('加载事件失败: $e');
    } finally {
      _loadingEvents = false;
      setState(() {});
    }
  }

  Future<void> _loadLastAction() async {
    try {
      _loadingAction = true;
      setState(() {});
      _lastAction = await _actionService.getLastAction();
    } catch (e) {
      debugPrint('加载行动失败: $e');
    } finally {
      _loadingAction = false;
      setState(() {});
    }
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: SingleChildScrollView(
          child: Column(
            children: [
              _buildHeader(),
              const SizedBox(height: 16),
              _buildAvatarStatusCard(context),
              const SizedBox(height: 16),
              _buildQuickActions(context),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildHeader() {
    return Padding(
      padding: const EdgeInsets.all(16),
      child: Row(
        mainAxisAlignment: MainAxisAlignment.spaceBetween,
        children: [
          const Text(
            'ME2',
            style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
          ),
          IconButton(
            icon: const Icon(Icons.notifications_outlined),
            onPressed: () {},
          ),
        ],
      ),
    );
  }

  Widget _buildAvatarStatusCard(BuildContext context) {
    final actionText = _lastAction != null
        ? '在${_lastAction!.sceneName}${_getActionText(_lastAction!.actionType)}'
        : '休息中';
    final icon = _lastAction != null
        ? _getIconForActionType(_lastAction!.actionType)
        : Icons.hotel;

    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        gradient: LinearGradient(
          colors: [Colors.purple.shade400, Colors.blue.shade400],
          begin: Alignment.topLeft,
          end: Alignment.bottomRight,
        ),
        borderRadius: BorderRadius.circular(20),
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '你的分身正在...',
            style: TextStyle(color: Colors.white70, fontSize: 14),
          ),
          const SizedBox(height: 8),
          _loadingAction
              ? const Text(
                  '加载中...',
                  style: TextStyle(
                    color: Colors.white,
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                  ),
                )
              : Text(
                  actionText,
                  style: const TextStyle(
                    color: Colors.white,
                    fontSize: 24,
                    fontWeight: FontWeight.bold,
                  ),
                ),
          const SizedBox(height: 16),
          Container(
            height: 150,
            decoration: BoxDecoration(
              color: Colors.white.withValues(alpha: 0.2),
              borderRadius: BorderRadius.circular(12),
            ),
            child: Center(
              child: Icon(icon, size: 64, color: Colors.white),
            ),
          ),
          const SizedBox(height: 16),
          ElevatedButton.icon(
            onPressed: () {},
            icon: const Icon(Icons.chat_bubble_outline),
            label: const Text('与 TA 对话'),
            style: ElevatedButton.styleFrom(
              backgroundColor: Colors.white,
              foregroundColor: Colors.purple,
              minimumSize: const Size(double.infinity, 48),
            ),
          ),
        ],
      ),
    );
  }

  String _getActionText(String actionType) {
    switch (actionType) {
      case 'exploration':
        return '探索';
      case 'social':
        return '社交';
      case 'study':
        return '学习';
      case 'creative':
        return '创作';
      case 'rest':
        return '休息';
      case 'play':
        return '娱乐';
      default:
        return '活动';
    }
  }

  IconData _getIconForActionType(String actionType) {
    switch (actionType) {
      case 'exploration':
        return Icons.explore;
      case 'social':
        return Icons.people;
      case 'study':
        return Icons.book;
      case 'creative':
        return Icons.palette;
      case 'rest':
        return Icons.spa;
      case 'play':
        return Icons.sports_esports;
      default:
        return Icons.event;
    }
  }

  Widget _buildQuickActions(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '最近动态',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 12),
          if (_loadingEvents)
            const Center(child: CircularProgressIndicator())
          else if (_recentEvents.isEmpty)
            const Padding(
              padding: EdgeInsets.all(16),
              child: Text('暂无事件'),
            )
          else
            _buildEventTimeline(),
        ],
      ),
    );
  }

  Widget _buildEventTimeline() {
    return Column(
      children: _recentEvents
          .asMap()
          .entries
          .map((entry) => _buildTimelineItem(
                entry.value,
                entry.key < _recentEvents.length - 1,
              ))
          .toList(),
    );
  }

  Widget _buildTimelineItem(Event event, bool showLine) {
    final color = _getColorForEventType(event.eventType);
    final icon = _getIconForEventType(event.eventType);
    final timeAgo = _formatTimeAgo(event.occurredAt);

    return IntrinsicHeight(
      child: Row(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Column(
            children: [
              Container(
                width: 40,
                height: 40,
                decoration: BoxDecoration(
                  color: color.withValues(alpha: 0.1),
                  shape: BoxShape.circle,
                ),
                child: Icon(icon, color: color, size: 20),
              ),
              if (showLine)
                Expanded(
                  child: Container(
                    width: 2,
                    margin: const EdgeInsets.symmetric(vertical: 4),
                    color: Colors.grey.shade300,
                  ),
                ),
            ],
          ),
          const SizedBox(width: 12),
          Expanded(
            child: Padding(
              padding: const EdgeInsets.only(bottom: 16),
              child: Column(
                crossAxisAlignment: CrossAxisAlignment.start,
                children: [
                  Text(
                    event.eventTitle,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                  const SizedBox(height: 2),
                  Text(
                    event.sceneName,
                    style: const TextStyle(
                      fontSize: 12,
                      color: Colors.grey,
                    ),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    event.eventText,
                    style: const TextStyle(color: Colors.grey),
                  ),
                  const SizedBox(height: 4),
                  Text(
                    timeAgo,
                    style: TextStyle(
                      fontSize: 12,
                      color: Colors.grey.shade600,
                    ),
                  ),
                ],
              ),
            ),
          ),
        ],
      ),
    );
  }

  Color _getColorForEventType(String eventType) {
    switch (eventType) {
      case 'exploration':
        return Colors.green;
      case 'social':
        return Colors.orange;
      case 'study':
        return Colors.blue;
      case 'creative':
        return Colors.purple;
      case 'rest':
        return Colors.teal;
      case 'play':
        return Colors.pink;
      default:
        return Colors.grey;
    }
  }

  IconData _getIconForEventType(String eventType) {
    switch (eventType) {
      case 'exploration':
        return Icons.explore;
      case 'social':
        return Icons.people;
      case 'study':
        return Icons.book;
      case 'creative':
        return Icons.palette;
      case 'rest':
        return Icons.spa;
      case 'play':
        return Icons.sports_esports;
      default:
        return Icons.event;
    }
  }

  String _formatTimeAgo(int timestamp) {
    final now = DateTime.now().millisecondsSinceEpoch ~/ 1000;
    final diff = now - timestamp;

    if (diff < 3600) {
      return '${diff ~/ 60}分钟前';
    } else if (diff < 86400) {
      return '${diff ~/ 3600}小时前';
    } else if (diff < 604800) {
      return '${diff ~/ 86400}天前';
    } else {
      return DateTime.fromMillisecondsSinceEpoch(timestamp * 1000)
          .toString()
          .substring(0, 10);
    }
  }
}
