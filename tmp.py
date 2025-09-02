from src.python_agent.session import SessionManager
manager = SessionManager()
session = manager.create_session()
session.add_message('user', 'Hello world')
manager.save_session(session)
print(f'Created session: {session.session_id}')

# List available sessions
sessions = manager.list_sessions()
print(f'Available sessions: {sessions}')