import { redirect } from 'react-router';

import session from '@/utils/session';

export default function authLoader() {
  const isAuthenticated = session.isAuthenticated();

  if (isAuthenticated) {
    return redirect('/');
  }

  return null;
}
