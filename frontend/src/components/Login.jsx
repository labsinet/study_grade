import React, { useState } from 'react';
import { API_BASE_URL } from '../config';

function Login({ setUser, setToken, setShowRegister }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  console.log('Login rendering', { username, password, error });

  const handleLogin = async (e) => {
    e.preventDefault();
    setError('');
    if (!username || !password) {
      setError('Заповніть усі поля');
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/api/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });

      if (response.ok) {
        const { user, token } = await response.json();
        setUser(user);
        setToken(token);
      } else {
        const errorText = await response.text();
        setError(errorText || 'Помилка входу');
      }
    } catch (err) {
      console.error('Login error:', err);
      setError('Не вдалося підключитися до сервера');
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="bg-white p-6 rounded shadow-md w-96">
        <h2 className="text-2xl mb-4">Вхід</h2>
        {error && <p className="text-red-500 mb-4">{error}</p>}
        <form onSubmit={handleLogin}>
          <input
            type="text"
            placeholder="Ім'я користувача"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
            className="w-full p-2 mb-4 border rounded"
            required
          />
          <input
            type="password"
            placeholder="Пароль"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full p-2 mb-4 border rounded"
            required
          />
          <button
            type="submit"
            className="w-full bg-blue-500 text-white p-2 rounded mb-2"
          >
            Увійти
          </button>
          <button
            type="button"
            onClick={() => setShowRegister(true)}
            className="w-full bg-gray-500 text-white p-2 rounded"
          >
            Реєстрація
          </button>
        </form>
      </div>
    </div>
  );
}

export default Login;