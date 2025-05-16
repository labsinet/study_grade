import React, { useState, Component } from 'react';
import Login from './components/Login';
import Register from './components/Register';
import Dashboard from './components/Dashboard';
import 'tailwindcss/tailwind.css';

class ErrorBoundary extends Component {
  state = { error: null };

  static getDerivedStateFromError(error) {
    return { error };
  }

  render() {
    if (this.state.error) {
      return (
        <div className="min-h-screen bg-gray-100 flex items-center justify-center">
          <div className="bg-white p-6 rounded shadow-md">
            <h1 className="text-2xl text-red-500">Помилка рендерингу</h1>
            <p>{this.state.error.message}</p>
          </div>
        </div>
      );
    }
    return this.props.children;
  }
}

function App() {
  const [user, setUser] = useState(null);
  const [token, setToken] = useState(null);
  const [showRegister, setShowRegister] = useState(false);

  console.log('App rendering', { user, token, showRegister });

  return (
    <ErrorBoundary>
      <div className="min-h-screen bg-gray-100">
        {user ? (
          <Dashboard user={user} token={token} setUser={setUser} setToken={setToken} />
        ) : showRegister ? (
          <Register setUser={setUser} setToken={setToken} setShowRegister={setShowRegister} />
        ) : (
          <Login setUser={setUser} setToken={setToken} setShowRegister={setShowRegister} />
        )}
      </div>
    </ErrorBoundary>
  );
}

export default App;