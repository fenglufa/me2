import 'package:flutter/material.dart';
import '../controllers/world_controller.dart';
import '../models/world_models.dart';
import '../../event/models/event_models.dart';
import '../../event/services/event_service.dart';

class WorldPage extends StatefulWidget {
  const WorldPage({super.key});

  @override
  State<WorldPage> createState() => _WorldPageState();
}

class _WorldPageState extends State<WorldPage> {
  final _controller = WorldController();
  final _eventService = EventService();
  List<Event> _recentEvents = [];
  bool _loadingEvents = false;

  @override
  void initState() {
    super.initState();
    _controller.addListener(() => setState(() {}));
    _controller.loadMaps();
    _loadRecentEvents();
  }

  Future<void> _loadRecentEvents() async {
    try {
      _loadingEvents = true;
      setState(() {});
      _recentEvents = await _eventService.getEventTimeline(pageSize: 5);
    } catch (e) {
      debugPrint('加载事件失败: $e');
    } finally {
      _loadingEvents = false;
      setState(() {});
    }
  }

  @override
  void dispose() {
    _controller.dispose();
    super.dispose();
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: _controller.isLoading && _controller.regions.isEmpty
            ? const Center(child: CircularProgressIndicator())
            : SingleChildScrollView(
                child: Column(
                  children: [
                    _buildHeader(),
                    const SizedBox(height: 16),
                    _buildRegionsSection(context),
                    const SizedBox(height: 16),
                    _buildRecentEventsSection(context),
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
            '第二空间',
            style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
          ),
          IconButton(
            icon: const Icon(Icons.map_outlined),
            onPressed: () {},
          ),
        ],
      ),
    );
  }

  Widget _buildRegionsSection(BuildContext context) {
    if (_controller.regions.isEmpty) {
      return const Padding(
        padding: EdgeInsets.all(16),
        child: Text('暂无区域数据'),
      );
    }

    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '世界区域',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 12),
          ..._controller.regions.map((region) => _buildRegionCard(region)),
        ],
      ),
    );
  }

  Widget _buildRegionCard(WorldRegion region) {
    final color = _getColorForRegion(region.atmosphere);
    final icon = _getIconForRegion(region.atmosphere);

    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: InkWell(
        onTap: () => _showRegionScenes(region),
        borderRadius: BorderRadius.circular(12),
        child: Container(
          padding: const EdgeInsets.all(16),
          decoration: BoxDecoration(
            borderRadius: BorderRadius.circular(12),
          ),
          child: Row(
            children: [
              Container(
                width: 60,
                height: 60,
                decoration: BoxDecoration(
                  color: color.withOpacity(0.1),
                  borderRadius: BorderRadius.circular(12),
                ),
                child: Icon(icon, color: color, size: 32),
              ),
              const SizedBox(width: 16),
              Expanded(
                child: Column(
                  crossAxisAlignment: CrossAxisAlignment.start,
                  children: [
                    Text(
                      region.name,
                      style: const TextStyle(
                        fontSize: 16,
                        fontWeight: FontWeight.bold,
                      ),
                    ),
                    const SizedBox(height: 4),
                    Text(
                      region.description,
                      style: const TextStyle(color: Colors.grey, fontSize: 12),
                      maxLines: 2,
                      overflow: TextOverflow.ellipsis,
                    ),
                    if (region.tags.isNotEmpty) ...[
                      const SizedBox(height: 4),
                      Wrap(
                        spacing: 4,
                        children: region.tags
                            .take(3)
                            .map((tag) => Chip(
                                  label: Text(tag),
                                  labelStyle: const TextStyle(fontSize: 10),
                                  padding: EdgeInsets.zero,
                                  materialTapTargetSize:
                                      MaterialTapTargetSize.shrinkWrap,
                                ))
                            .toList(),
                      ),
                    ],
                  ],
                ),
              ),
              const Icon(Icons.chevron_right, color: Colors.grey),
            ],
          ),
        ),
      ),
    );
  }

  Color _getColorForRegion(String atmosphere) {
    switch (atmosphere) {
      case '繁华':
        return Colors.blue;
      case '宁静':
        return Colors.green;
      case '学术':
        return Colors.orange;
      case '艺术':
        return Colors.purple;
      default:
        return Colors.grey;
    }
  }

  IconData _getIconForRegion(String atmosphere) {
    switch (atmosphere) {
      case '繁华':
        return Icons.location_city;
      case '宁静':
        return Icons.nature;
      case '学术':
        return Icons.school;
      case '艺术':
        return Icons.palette;
      default:
        return Icons.place;
    }
  }

  void _showRegionScenes(WorldRegion region) async {
    await _controller.loadScenes(region.id);
    if (!mounted) return;

    showModalBottomSheet(
      context: context,
      isScrollControlled: true,
      builder: (context) => DraggableScrollableSheet(
        initialChildSize: 0.7,
        minChildSize: 0.5,
        maxChildSize: 0.95,
        expand: false,
        builder: (context, scrollController) => Container(
          padding: const EdgeInsets.all(16),
          child: Column(
            crossAxisAlignment: CrossAxisAlignment.start,
            children: [
              Center(
                child: Container(
                  width: 40,
                  height: 4,
                  margin: const EdgeInsets.only(bottom: 16),
                  decoration: BoxDecoration(
                    color: Colors.grey.shade300,
                    borderRadius: BorderRadius.circular(2),
                  ),
                ),
              ),
              Text(
                region.name,
                style:
                    const TextStyle(fontSize: 20, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 8),
              Text(
                region.description,
                style: const TextStyle(color: Colors.grey),
              ),
              const SizedBox(height: 16),
              const Text(
                '场景列表',
                style: TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
              ),
              const SizedBox(height: 12),
              Expanded(
                child: _controller.isLoading
                    ? const Center(child: CircularProgressIndicator())
                    : _controller.scenes.isEmpty
                        ? const Center(child: Text('暂无场景'))
                        : ListView.builder(
                            controller: scrollController,
                            itemCount: _controller.scenes.length,
                            itemBuilder: (context, index) =>
                                _buildSceneCard(_controller.scenes[index]),
                          ),
              ),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildSceneCard(WorldScene scene) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(12),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Row(
              children: [
                Expanded(
                  child: Text(
                    scene.name,
                    style: const TextStyle(
                      fontSize: 16,
                      fontWeight: FontWeight.bold,
                    ),
                  ),
                ),
                Container(
                  padding:
                      const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
                  decoration: BoxDecoration(
                    color: Colors.blue.shade100,
                    borderRadius: BorderRadius.circular(12),
                  ),
                  child: Text(
                    scene.atmosphere,
                    style: TextStyle(
                      fontSize: 10,
                      color: Colors.blue.shade700,
                    ),
                  ),
                ),
              ],
            ),
            const SizedBox(height: 8),
            Text(
              scene.description,
              style: const TextStyle(color: Colors.grey, fontSize: 12),
            ),
            const SizedBox(height: 8),
            Wrap(
              spacing: 8,
              runSpacing: 4,
              children: [
                if (scene.features.hasWifi)
                  _buildFeatureChip(Icons.wifi, 'WiFi'),
                if (scene.features.hasFood)
                  _buildFeatureChip(Icons.restaurant, '餐饮'),
                if (scene.features.isQuiet)
                  _buildFeatureChip(Icons.volume_off, '安静'),
                if (scene.features.isIndoor)
                  _buildFeatureChip(Icons.home, '室内'),
              ],
            ),
          ],
        ),
      ),
    );
  }

  Widget _buildFeatureChip(IconData icon, String label) {
    return Container(
      padding: const EdgeInsets.symmetric(horizontal: 8, vertical: 4),
      decoration: BoxDecoration(
        color: Colors.grey.shade200,
        borderRadius: BorderRadius.circular(12),
      ),
      child: Row(
        mainAxisSize: MainAxisSize.min,
        children: [
          Icon(icon, size: 12, color: Colors.grey.shade700),
          const SizedBox(width: 4),
          Text(
            label,
            style: TextStyle(fontSize: 10, color: Colors.grey.shade700),
          ),
        ],
      ),
    );
  }

  Widget _buildRecentEventsSection(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                '最近事件',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              TextButton(
                onPressed: () {},
                child: const Text('查看全部'),
              ),
            ],
          ),
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
