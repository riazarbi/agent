from python_agent.session import SessionManager
sm = SessionManager()
sessions = sm.list_sessions()
print(f'Found {len(sessions)} sessions')
for sid in sessions[-3:]:  # Show last 3
    print(f'  {sid}')