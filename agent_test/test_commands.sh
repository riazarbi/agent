#!/bin/bash

echo "=== Manual Testing Guide for Unified Prompt Handling ==="
echo
echo "Test files created in: $(pwd)/agent_test"
echo "Agent binary location: $(pwd)/agent"
echo
echo "Run these commands to test different scenarios:"
echo
echo "1. Test new --prompt flag with YAML (complete poem):"
echo "   ./agent --prompt agent_test/complete_prompt.yml"
echo
echo "2. Test new --prompt flag with Markdown:"
echo "   ./agent --prompt agent_test/poem_intro.md"
echo
echo "3. Test new --preprompt flag with legacy format:"
echo "   ./agent --preprompt agent_test/legacy_preprompts"
echo
echo "4. Test new --preprompt flag with YAML + interactive:"
echo "   ./agent --preprompt agent_test/complete_prompt.yml"
echo
echo "5. Test simple YAML example (demonstrates spec features):"
echo "   ./agent --prompt agent_test/simple_example.yml"
echo
echo "Expected result: The agent should have the complete poem loaded and"
echo "be able to answer questions about it, like:"
echo "- 'What poem did you just recite?'"
echo "- 'Who wrote that poem?'"
echo "- 'What are the last two lines?'"
echo "- 'How many stanzas were there?'"
echo
echo "Files in test directory:"
ls -la agent_test/