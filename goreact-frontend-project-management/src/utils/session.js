const SESSION_KEY = 'my-session';

const session = {
  setSession(user) {
    localStorage.setItem(SESSION_KEY, JSON.stringify(user));
  },
  getSession() {
    const user = JSON.parse(localStorage.getItem(SESSION_KEY));
    return user;
  },
  getToken() {
    const session = this.getSession();
    return session?.access_token ?? null;
  },
  clearSession() {
    localStorage.removeItem(SESSION_KEY);
  },
  isAuthenticated() {
    const session = localStorage.getItem(SESSION_KEY);
    return !!session;
  },
};

export default session;
