import { redirect } from 'react-router';

import session from '@/utils/session';

export default function sidebarLoader() {
  const isAuthenticated = session.isAuthenticated();

  if (!isAuthenticated) {
    return redirect('/login');
  }

  return null;
}
