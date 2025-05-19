import type { FC, FormEvent } from 'react';
import { useAuth } from 'react-oidc-context';
import { useEffect, useState } from 'react';
import { Message } from '../Message.tsx';

interface MessageData {
  author: string;
  content: string;
}

export const Home: FC = () => {
  const auth = useAuth();
  const [messages, setMessages] = useState<MessageData[]>([]);
  const [error, setError] = useState<string | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [newMessage, setNewMessage] = useState<string>(''); // State for the new message

  useEffect(() => {
    const fetchMessages = async () => {
      try {
        setLoading(true);
        setError(null);

        const response = await fetch('http://localhost:8080/messages?limit=100&offset=0', {
          headers: {
            Authorization: `Bearer ${auth.user?.access_token}`,
          },
        });

        if (!response.ok) {
          throw new Error(`Failed to fetch messages: ${response.statusText}`);
        }

        const data = await response.json();
        setMessages(data); // Assuming the backend returns an array of messages
      } catch (err: any) {
        setError(err.message);
      } finally {
        setLoading(false);
      }
    };

    fetchMessages();
  }, [auth.user?.access_token]);

  const handleSendMessage = async (e: FormEvent) => {
    e.preventDefault(); // Prevent form submission from reloading the page

    try {
      const response = await fetch('http://localhost:8080/messages', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          Authorization: `Bearer ${auth.user?.access_token}`,
        },
        body: JSON.stringify({ content: newMessage }),
      });

      if (!response.ok) {
        throw new Error(`Failed to send message: ${response.statusText}`);
      }

      const createdMessage = await response.json();
      setMessages((prevMessages) => [createdMessage, ...prevMessages]); // Add the new message to the list
      setNewMessage(''); // Clear the textbox
    } catch (err: any) {
      setError(err.message);
    }
  };

  if (loading) {
    return <p>Loading messages...</p>;
  }

  if (error) {
    return <p>Error: {error}</p>;
  }

  return (
    <>
      <h1>Home</h1>
      <p>Inspecting the result of the <code>useAuth()</code> hook.</p>

      <h2>Messages</h2>
      {messages.length > 0 ? (
        messages.map((message, index) => (
          <Message key={index} author={message.author} content={message.content} />
        ))
      ) : (
        <p>No messages found.</p>
      )}

      <h2>Send a Message</h2>
      <form onSubmit={handleSendMessage}>
        <textarea
          value={newMessage}
          onChange={(e) => setNewMessage(e.target.value)}
          placeholder="Write your message here..."
          rows={4}
          style={{ width: '100%', marginBottom: '8px' }}
        />
        <button type="submit" disabled={!newMessage.trim()}>
          Send
        </button>
      </form>
    </>
  );
};