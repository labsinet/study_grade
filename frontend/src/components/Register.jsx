import React, { useState } from 'react';
import { API_BASE_URL } from '../config';

function Register({ setUser, setToken, setShowRegister }) {
  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  console.log('Register rendering', { username, password, error });

  const handleRegister = async (e) => {
    e.preventDefault();
    setError('');
    console.log('Submitting register with POST to', `${API_BASE_URL}/api/register`);
    if (!username || !password) {
      setError('Заповніть усі поля');
      return;
    }
    if (username.length < 3 || username.length > 50) {
      setError('Ім’я користувача має бути від 3 до 50 символів');
      return;
    }
    if (password.length < 8) {
      setError('Пароль має бути щонайменше 8 символів');
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/api/register`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ username, password }),
      });
      console.log('Register response status:', response.status);

      if (response.ok) {
        const user = await response.json();
        const loginResponse = await fetch(`${API_BASE_URL}/api/login`, {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({ username, password }),
        });
        if (loginResponse.ok) {
          const { user: newUser, token } = await loginResponse.json();
          setUser(newUser);
          setToken(token);
        } else {
          setError('Не вдалося увійти після реєстрації');
        }
      } else {
        const errorText = await response.text();
        setError(errorText || 'Помилка реєстрації');
      }
    } catch (err) {
      console.error('Register fetch error:', err);
      setError('Не вдалося підключитися до сервера');
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="bg-white p-6 rounded shadow-md w-96">
        <h2 className="text-2xl mb-4">Реєстрація</h2>
        {error && <p className="text-red-500 mb-4">{error}</p>}
        <form onSubmit={handleRegister}>
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
            disabled={!username || !password}
          >
            Зареєструватися
          </button>
          <button
            type="button"
            onClick={() => setShowRegister(false)}
            className="w-full bg-gray-500 text-white p-2 rounded"
          >
            Повернутися до входу
          </button>
        </form>
      </div>
    </div>
  );
}

export default Register;