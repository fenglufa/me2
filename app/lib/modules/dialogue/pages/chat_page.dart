import 'dart:convert';
import 'package:flutter/material.dart';
import 'package:web_socket_channel/web_socket_channel.dart';
import '../services/dialogue_service.dart';
import '../models/dialogue_models.dart';

class ChatPage extends StatefulWidget {
  final int avatarId;
  final String avatarName;

  const ChatPage({
    super.key,
    required this.avatarId,
    required this.avatarName,
  });

  @override
  State<ChatPage> createState() => _ChatPageState();
}

class _ChatPageState extends State<ChatPage> {
  final _dialogueService = DialogueService();
  final _messageController = TextEditingController();
  final _scrollController = ScrollController();
  final List<Message> _messages = [];
  WebSocketChannel? _channel;
  int? _sessionId;
  bool _loading = false;
  bool _streaming = false;
  String _streamingContent = '';

  @override
  void initState() {
    super.initState();
    _initSession();
  }

  @override
  void dispose() {
    _channel?.sink.close();
    _messageController.dispose();
    _scrollController.dispose();
    super.dispose();
  }

  Future<void> _initSession() async {
    setState(() => _loading = true);
    try {
      final sessions = await _dialogueService.getSessions(widget.avatarId, pageSize: 1);
      if (sessions.isNotEmpty) {
        _sessionId = sessions.first.id;
        await _loadMessages();
      } else {
        _sessionId = await _dialogueService.createSession(widget.avatarId);
      }
      _channel = await _dialogueService.connectStream(widget.avatarId);
    } catch (e) {
      if (mounted) {
        ScaffoldMessenger.of(context).showSnackBar(
          SnackBar(content: Text('初始化失败: $e')),
        );
      }
    } finally {
      setState(() => _loading = false);
    }
  }

  Future<void> _loadMessages() async {
    if (_sessionId == null) return;
    try {
      final messages = await _dialogueService.getMessages(_sessionId!);
      setState(() => _messages.addAll(messages));
      _scrollToBottom();
    } catch (e) {
      debugPrint('加载消息失败: $e');
    }
  }

  void _sendMessage() {
    if (_messageController.text.trim().isEmpty || _sessionId == null || _streaming) return;

    final userMessage = _messageController.text.trim();
    _messageController.clear();

    setState(() {
      _messages.add(Message(
        id: 0,
        sessionId: _sessionId!,
        role: 'user',
        content: userMessage,
        createdAt: DateTime.now().millisecondsSinceEpoch ~/ 1000,
      ));
      _streaming = true;
      _streamingContent = '';
    });

    _scrollToBottom();

    _channel?.sink.add(jsonEncode({
      'session_id': _sessionId,
      'content': userMessage,
    }));

    _channel?.stream.listen(
      (data) {
        final response = jsonDecode(data);
        setState(() {
          if (response['content'] != null && response['content'].isNotEmpty) {
            _streamingContent += response['content'];
          }
          if (response['done'] == true) {
            if (_streamingContent.isNotEmpty) {
              _messages.add(Message(
                id: response['message_id'] ?? 0,
                sessionId: _sessionId!,
                role: 'assistant',
                content: _streamingContent,
                createdAt: DateTime.now().millisecondsSinceEpoch ~/ 1000,
              ));
            }
            _streaming = false;
            _streamingContent = '';
          }
        });
        _scrollToBottom();
      },
      onError: (error) {
        setState(() {
          _streaming = false;
          _streamingContent = '';
        });
        if (mounted) {
          ScaffoldMessenger.of(context).showSnackBar(
            SnackBar(content: Text('发送失败: $error')),
          );
        }
      },
    );
  }

  void _scrollToBottom() {
    WidgetsBinding.instance.addPostFrameCallback((_) {
      if (_scrollController.hasClients) {
        _scrollController.animateTo(
          _scrollController.position.maxScrollExtent,
          duration: const Duration(milliseconds: 300),
          curve: Curves.easeOut,
        );
      }
    });
  }

  @override
  Widget build(BuildContext context) {
    return Scaffold(
      appBar: AppBar(
        title: Text(widget.avatarName),
      ),
      body: _loading
          ? const Center(child: CircularProgressIndicator())
          : Column(
              children: [
                Expanded(
                  child: ListView.builder(
                    controller: _scrollController,
                    padding: const EdgeInsets.all(16),
                    itemCount: _messages.length + (_streaming ? 1 : 0),
                    itemBuilder: (context, index) {
                      if (index == _messages.length) {
                        return _buildMessageBubble('assistant', _streamingContent);
                      }
                      final message = _messages[index];
                      return _buildMessageBubble(message.role, message.content);
                    },
                  ),
                ),
                _buildInputBar(),
              ],
            ),
    );
  }

  Widget _buildMessageBubble(String role, String content) {
    final isUser = role == 'user';
    return Align(
      alignment: isUser ? Alignment.centerRight : Alignment.centerLeft,
      child: Container(
        margin: const EdgeInsets.only(bottom: 12),
        padding: const EdgeInsets.symmetric(horizontal: 16, vertical: 12),
        constraints: BoxConstraints(maxWidth: MediaQuery.of(context).size.width * 0.75),
        decoration: BoxDecoration(
          color: isUser ? Colors.blue : Colors.grey.shade200,
          borderRadius: BorderRadius.circular(16),
        ),
        child: Text(
          content,
          style: TextStyle(color: isUser ? Colors.white : Colors.black87),
        ),
      ),
    );
  }

  Widget _buildInputBar() {
    return Container(
      padding: const EdgeInsets.all(16),
      decoration: BoxDecoration(
        color: Colors.white,
        boxShadow: [
          BoxShadow(
            color: Colors.black.withValues(alpha: 0.05),
            blurRadius: 10,
          ),
        ],
      ),
      child: Row(
        children: [
          Expanded(
            child: TextField(
              controller: _messageController,
              enabled: !_streaming,
              decoration: InputDecoration(
                hintText: '输入消息...',
                border: OutlineInputBorder(
                  borderRadius: BorderRadius.circular(24),
                  borderSide: BorderSide.none,
                ),
                filled: true,
                fillColor: Colors.grey.shade100,
                contentPadding: const EdgeInsets.symmetric(horizontal: 20, vertical: 12),
              ),
              onSubmitted: (_) => _sendMessage(),
            ),
          ),
          const SizedBox(width: 8),
          IconButton(
            onPressed: _streaming ? null : _sendMessage,
            icon: Icon(_streaming ? Icons.hourglass_empty : Icons.send),
            color: Colors.blue,
          ),
        ],
      ),
    );
  }
}
