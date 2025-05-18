import type { FC } from 'react';
import { useAuth } from 'react-oidc-context';
import { Message } from '../Message.tsx';

interface Row {
  label: string;
  value: React.ReactNode;
}

const createRows = (data?: unknown): Row[] => {
  if (!data) {
    return [];
  }

  return Object.entries(data).map(([key, value]) => {
    return {
      label: key,
      value: JSON.stringify(value),
    };
  });
};

export const Home: FC = () => {
  const auth = useAuth();

  return (
    <>
      <h1>Home</h1>
      <p>
        Inspecting the result of the <code>useAuth()</code> hook.
      </p>

      <h2>
        <code>auth.user?.profile</code>
      </h2>
      <Message author={auth.user?.profile.name} content={"a"} />
    </>
  );
};
