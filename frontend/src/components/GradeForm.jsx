import React, { useState } from 'react';
import { API_BASE_URL } from '../config';

function GradeForm({ userId, token }) {
  const [formData, setFormData] = useState({
    date: '',
    semester: '',
    subject: '',
    group: '',
    totalStudents: '',
    grade5: '',
    grade4: '',
    grade3: '',
    grade2: '',
    notPassed: '',
  });
  const [error, setError] = useState('');

  console.log('GradeForm rendering', { userId, formData, error });

  const handleChange = (e) => {
    setFormData({ ...formData, [e.target.name]: e.target.value });
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setError('');

    if (!formData.date || !formData.semester || !formData.subject || !formData.group ||
        !formData.totalStudents || formData.grade5 === '' || formData.grade4 === '' ||
        formData.grade3 === '' || formData.grade2 === '' || formData.notPassed === '') {
      setError('Заповніть усі поля');
      return;
    }

    const totalStudents = parseInt(formData.totalStudents);
    const grade5 = parseInt(formData.grade5);
    const grade4 = parseInt(formData.grade4);
    const grade3 = parseInt(formData.grade3);
    const grade2 = parseInt(formData.grade2);
    const notPassed = parseInt(formData.notPassed);

    if (isNaN(totalStudents) || totalStudents < 1) {
      setError('Кількість студентів має бути більше 0');
      return;
    }
    if (grade5 < 0 || grade4 < 0 || grade3 < 0 || grade2 < 0 || notPassed < 0) {
      setError('Оцінки не можуть бути від’ємними');
      return;
    }
    if (grade5 + grade4 + grade3 + grade2 + notPassed !== totalStudents) {
      setError('Сума оцінок має дорівнювати загальній кількості студентів');
      return;
    }

    try {
      const response = await fetch(`${API_BASE_URL}/api/grades`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`,
        },
        body: JSON.stringify({
          date: formData.date,
          semester: parseInt(formData.semester),
          subject: formData.subject,
          group: formData.group,
          total_students: totalStudents,
          grade_5: grade5,
          grade_4: grade4,
          grade_3: grade3,
          grade_2: grade2,
          not_passed: notPassed,
          user_id: userId,
        }),
      });

      if (response.ok) {
        alert('Дані збережено');
        setFormData({
          date: '',
          semester: '',
          subject: '',
          group: '',
          totalStudents: '',
          grade5: '',
          grade4: '',
          grade3: '',
          grade2: '',
          notPassed: '',
        });
      } else {
        const errorText = await response.text();
        setError(errorText || 'Помилка збереження');
      }
    } catch (err) {
      console.error('GradeForm error:', err);
      setError('Не вдалося підключитися до сервера');
    }
  };

  return (
    <div className="bg-white p-6 rounded shadow-md mb-6">
      <h2 className="text-2xl mb-4">Додати оцінки</h2>
      {error && <p className="text-red-500 mb-4">{error}</p>}
      <form onSubmit={handleSubmit}>
        <input
          type="date"
          name="date"
          value={formData.date}
          onChange={handleChange}
          className="w-full p-2 mb-2 border rounded"
          required
        />
        <input
          type="number"
          name="semester"
          placeholder="Семестр"
          value={formData.semester}
          onChange={handleChange}
          className="w-full p-2 mb-2 border rounded"
          required
        />
        <input
          type="text"
          name="subject"
          placeholder="Предмет"
          value={formData.subject}
          onChange={handleChange}
          className="w-full p-2 mb-2 border rounded"
          required
        />
        <input
          type="text"
          name="group"
          placeholder="Група"
          value={formData.group}
          onChange={handleChange}
          className="w-full p-2 mb-2 border rounded"
          required
        />
        <input
          type="number"
          name="totalStudents"
          placeholder="Кількість студентів"
          value={formData.totalStudents}
          onChange={handleChange}
          className="w-full p-2 mb-2 border rounded"
          required
        />
        <input
          type="number"
          name="grade5"
          placeholder="Оцінка 5"
          value={formData.grade5}
          onChange={handleChange}
          className="w-full p-2 mb-2 border rounded"
          required
        />
        <input
          type="number"
          name="grade4"
          placeholder="Оцінка 4"
          value={formData.grade4}
          onChange={handleChange}
          className="w-full p-2 mb-2 border rounded"
          required
        />
        <input
          type="number"
          name="grade3"
          placeholder="Оцінка 3"
          value={formData.grade3}
          onChange={handleChange}
          className="w-full p-2 mb-2 border rounded"
          required
        />
        <input
          type="number"
          name="grade2"
          placeholder="Оцінка 2"
          value={formData.grade2}
          onChange={handleChange}
          className="w-full p-2 mb-2 border rounded"
          required
        />
        <input
          type="number"
          name="notPassed"
          placeholder="Не атестовані"
          value={formData.notPassed}
          onChange={handleChange}
          className="w-full p-2 mb-2 border rounded"
          required
        />
        <button
          type="submit"
          className="w-full bg-blue-500 text-white p-2 rounded"
        >
          Зберегти
        </button>
      </form>
    </div>
  );
}

export default GradeForm;