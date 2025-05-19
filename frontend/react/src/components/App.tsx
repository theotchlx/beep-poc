import type { FC } from 'react';
import { Route, Routes } from 'react-router-dom';
import { appRoutes } from '../constants.ts';
import { Home } from './routes/Home.tsx';
import { NotFound } from './routes/NotFound.tsx';
import { Playground } from './routes/Playground/Playground.tsx';

const originalFetch = window.fetch;
window.fetch = async (...args) => {
  console.log('Fetch request:', args);
  const response = await originalFetch(...args);
  console.log('Fetch response:', response);
  return response;
};

export const App: FC = () => {
  return (
    <Routes>
      <Route path={appRoutes.home}>
        <Route index={true} element={<Home />} />
        <Route path={appRoutes.notFound} element={<NotFound />} />
        <Route path={appRoutes.playground} element={<Playground />} />
      </Route>
    </Routes>
  );
};
