import type { FC } from 'react';
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
    </>
  );
};