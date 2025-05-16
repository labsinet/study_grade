import React from 'react';
import GradeForm from './GradeForm';
import GradeTable from './GradeTable';

function Dashboard({ user, token, setUser, setToken }) {
  console.log('Dashboard rendering', { user, token });

  return (
    <div className="container mx-auto p-4">
      <div className="flex justify-between items-center mb-4">
        <h1 className="text-3xl">Вітаємо, {user.username}!</h1>
        <button
          onClick={() => {
            setUser(null);
            setToken(null);
          }}
          className="bg-red-500 text-white px-4 py-2 rounded"
        >
          Вийти
        </button>
      </div>
      <GradeForm userId={user.id} token={token} />
      <GradeTable userId={user.id} token={token} />
    </div>
  );
}

export default Dashboard;