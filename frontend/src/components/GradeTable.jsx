import React, { useState, useEffect } from 'react';
import { API_BASE_URL } from '../config';

function GradeTable({ userId, token }) {
  const [grades, setGrades] = useState([]);
  const [error, setError] = useState('');

  console.log('GradeTable rendering', { userId, token, grades });

  useEffect(() => {
    const fetchGrades = async () => {
      try {
        const response = await fetch(`${API_BASE_URL}/api/grades`, {
          headers: {
            'Authorization': `Bearer ${token}`,
          },
        });
        if (response.ok) {
          const data = await response.json();
          setGrades(data);
        } else {
          const errorText = await response.text();
          setError(errorText || 'Помилка завантаження оцінок');
        }
      } catch (err) {
        console.error('GradeTable error:', err);
        setError('Не вдалося підключитися до сервера');
      }
    };
    fetchGrades();
  }, [token]);

  const calculateAverages = () => {
    if (grades.length === 0) return null;
    const totals = grades.reduce(
      (acc, grade) => ({
        totalStudents: acc.totalStudents + grade.total_students,
        grade5: acc.grade5 + grade.grade_5,
        grade4: acc.grade4 + grade.grade_4,
        grade3: acc.grade3 + grade.grade_3,
        grade2: acc.grade2 + grade.grade_2,
        notPassed: acc.notPassed + grade.not_passed,
        averageScore: acc.averageScore + grade.average_score * grade.total_students,
      }),
      {
        totalStudents: 0,
        grade5: 0,
        grade4: 0,
        grade3: 0,
        grade2: 0,
        notPassed: 0,
        averageScore: 0,
      }
    );
    return {
      avgScore: (totals.averageScore / totals.totalStudents).toFixed(2),
      successRate: ((totals.grade5 + totals.grade4 + totals.grade3) / totals.totalStudents * 100).toFixed(2),
      qualityRate: ((totals.grade5 + totals.grade4) / totals.totalStudents * 100).toFixed(2),
    };
  };

  const averages = calculateAverages();

  return (
    <div className="bg-white p-6 rounded shadow-md">
      <h2 className="text-2xl mb-4">Оцінки</h2>
      {error && <p className="text-red-500 mb-4">{error}</p>}
      {grades.length === 0 && !error && <p>Оцінки відсутні</p>}
      {grades.length > 0 && (
        <table className="w-full border-collapse">
          <thead>
            <tr className="bg-gray-200">
              <th className="border p-2">Дата</th>
              <th className="border p-2">Семестр</th>
              <th className="border p-2">Предмет</th>
              <th className="border p-2">Група</th>
              <th className="border p-2">Студенти</th>
              <th className="border p-2">5</th>
              <th className="border p-2">4</th>
              <th className="border p-2">3</th>
              <th className="border p-2">2</th>
              <th className="border p-2">Не атест.</th>
              <th className="border p-2">Середній бал</th>
              <th className="border p-2">Успішність (%)</th>
              <th className="border p-2">Якість (%)</th>
            </tr>
          </thead>
          <tbody>
            {grades.map((grade) => (
              <tr key={grade.id}>
                <td className="border p-2">{new Date(grade.date).toLocaleDateString()}</td>
                <td className="border p-2">{grade.semester}</td>
                <td className="border p-2">{grade.subject}</td>
                <td className="border p-2">{grade.group}</td>
                <td className="border p-2">{grade.total_students}</td>
                <td className="border p-2">{grade.grade_5}</td>
                <td className="border p-2">{grade.grade_4}</td>
                <td className="border p-2">{grade.grade_3}</td>
                <td className="border p-2">{grade.grade_2}</td>
                <td className="border p-2">{grade.not_passed}</td>
                <td className="border p-2">{grade.average_score.toFixed(2)}</td>
                <td className="border p-2">{grade.success_rate.toFixed(2)}</td>
                <td className="border p-2">{grade.quality_rate.toFixed(2)}</td>
              </tr>
            ))}
            {averages && (
              <tr className="bg-gray-100 font-bold">
                <td className="border p-2" colSpan="4">Середнє за семестр</td>
                <td className="border p-2"></td>
                <td className="border p-2"></td>
                <td className="border p-2"></td>
                <td className="border p-2"></td>
                <td className="border p-2"></td>
                <td className="border p-2"></td>
                <td className="border p-2">{averages.avgScore}</td>
                <td className="border p-2">{averages.successRate}</td>
                <td className="border p-2">{averages.qualityRate}</td>
              </tr>
            )}
          </tbody>
        </table>
      )}
    </div>
  );
}

export default GradeTable;