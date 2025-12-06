import 'package:flutter/material.dart';

class AvatarPage extends StatelessWidget {
  const AvatarPage({super.key});

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      body: SafeArea(
        child: SingleChildScrollView(
          child: Column(
            children: [
              _buildHeader(context),
              _buildPersonalityCard(context),
              const SizedBox(height: 16),
              _buildDiarySection(context),
              const SizedBox(height: 16),
              _buildGrowthSection(context),
            ],
          ),
        ),
      ),
    );
  }

  Widget _buildHeader(BuildContext context) {
    return Container(
      padding: const EdgeInsets.all(20),
      child: Row(
        children: [
          CircleAvatar(
            radius: 40,
            backgroundColor: Colors.purple.shade100,
            child: const Icon(Icons.person, size: 40),
          ),
          const SizedBox(width: 16),
          const Expanded(
            child: Column(
              crossAxisAlignment: CrossAxisAlignment.start,
              children: [
                Text(
                  '我的分身',
                  style: TextStyle(fontSize: 24, fontWeight: FontWeight.bold),
                ),
                SizedBox(height: 4),
                Text('成长等级 Lv.5', style: TextStyle(color: Colors.grey)),
              ],
            ),
          ),
          IconButton(
            icon: const Icon(Icons.edit_outlined),
            onPressed: () {},
          ),
        ],
      ),
    );
  }

  Widget _buildPersonalityCard(BuildContext context) {
    return Container(
      margin: const EdgeInsets.symmetric(horizontal: 16),
      padding: const EdgeInsets.all(20),
      decoration: BoxDecoration(
        color: Colors.white,
        borderRadius: BorderRadius.circular(16),
        boxShadow: [
          BoxShadow(
            color: Colors.black.withOpacity(0.05),
            blurRadius: 10,
            offset: const Offset(0, 2),
          ),
        ],
      ),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '性格面板',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 16),
          Container(
            height: 200,
            decoration: BoxDecoration(
              color: Colors.purple.shade50,
              borderRadius: BorderRadius.circular(12),
            ),
            child: const Center(
              child: Text('6维雷达图'),
            ),
          ),
          const SizedBox(height: 16),
          _buildPersonalityItem('温暖度', 0.8, Colors.red),
          _buildPersonalityItem('冒险性', 0.6, Colors.orange),
          _buildPersonalityItem('社交性', 0.7, Colors.blue),
          _buildPersonalityItem('创造力', 0.9, Colors.purple),
          _buildPersonalityItem('平静度', 0.5, Colors.green),
          _buildPersonalityItem('活力值', 0.8, Colors.pink),
        ],
      ),
    );
  }

  Widget _buildPersonalityItem(String label, double value, Color color) {
    return Padding(
      padding: const EdgeInsets.only(bottom: 8),
      child: Row(
        children: [
          SizedBox(
            width: 60,
            child: Text(label, style: const TextStyle(fontSize: 12)),
          ),
          Expanded(
            child: LinearProgressIndicator(
              value: value,
              backgroundColor: Colors.grey.shade200,
              valueColor: AlwaysStoppedAnimation(color),
              minHeight: 8,
              borderRadius: BorderRadius.circular(4),
            ),
          ),
          const SizedBox(width: 8),
          Text('${(value * 100).toInt()}', style: const TextStyle(fontSize: 12)),
        ],
      ),
    );
  }

  Widget _buildDiarySection(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          Row(
            mainAxisAlignment: MainAxisAlignment.spaceBetween,
            children: [
              const Text(
                '分身日记',
                style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
              ),
              TextButton(
                onPressed: () {},
                child: const Text('查看全部'),
              ),
            ],
          ),
          _buildDiaryCard('今天去了森林', '心情很好，遇见了小鹿', '2小时前'),
          _buildDiaryCard('读完了一本书', '学到了很多新知识', '昨天'),
        ],
      ),
    );
  }

  Widget _buildDiaryCard(String title, String content, String time) {
    return Card(
      margin: const EdgeInsets.only(bottom: 12),
      child: Padding(
        padding: const EdgeInsets.all(16),
        child: Column(
          crossAxisAlignment: CrossAxisAlignment.start,
          children: [
            Text(
              title,
              style: const TextStyle(fontSize: 16, fontWeight: FontWeight.bold),
            ),
            const SizedBox(height: 8),
            Text(content, style: const TextStyle(color: Colors.grey)),
            const SizedBox(height: 8),
            Text(time, style: TextStyle(fontSize: 12, color: Colors.grey.shade600)),
          ],
        ),
      ),
    );
  }

  Widget _buildGrowthSection(BuildContext context) {
    return Padding(
      padding: const EdgeInsets.symmetric(horizontal: 16),
      child: Column(
        crossAxisAlignment: CrossAxisAlignment.start,
        children: [
          const Text(
            '成长节点',
            style: TextStyle(fontSize: 18, fontWeight: FontWeight.bold),
          ),
          const SizedBox(height: 12),
          Wrap(
            spacing: 8,
            runSpacing: 8,
            children: [
              _buildMilestoneChip('第一次旅行', true),
              _buildMilestoneChip('第一次社交', true),
              _buildMilestoneChip('第一次创作', false),
              _buildMilestoneChip('第一次探索', true),
            ],
          ),
        ],
      ),
    );
  }

  Widget _buildMilestoneChip(String label, bool achieved) {
    return Chip(
      label: Text(label),
      avatar: Icon(
        achieved ? Icons.check_circle : Icons.circle_outlined,
        size: 16,
        color: achieved ? Colors.green : Colors.grey,
      ),
      backgroundColor: achieved ? Colors.green.shade50 : Colors.grey.shade100,
    );
  }
}
