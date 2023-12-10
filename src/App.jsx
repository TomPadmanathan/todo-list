// React
import { useState, useEffect } from 'react';

// CSS
import './App.css';

// Assets
import { HiOutlineTrash, HiPlus } from 'react-icons/hi';

export default function App() {
	const [todos, setTodos] = useState([]);
	const [error, setError] = useState('');
	const [newTodo, setNewTodo] = useState('');

	function getDateTimeFromTimeStamp(timestamp) {
		const date = new Date(timestamp * 1000);

		const dateStr = `${date.getHours().toString().padStart(2, '0')}:${date
			.getMinutes()
			.toString()
			.padStart(2, '0')}`;
		const timeStr = `${date.getDate().toString().padStart(2, '0')}/${(
			date.getMonth() + 1
		)
			.toString()
			.padStart(2, '0')}/${date.getFullYear()}`;

		return dateStr + ' • ' + timeStr;
	}

	useEffect(() => {
		async function fetchTodos() {
			try {
				const response = await fetch('http://localhost:3001/todos', {
					method: 'GET',
					headers: { 'Content-Type': 'application/json' },
				});
				if (!response.ok) {
					switch (response.status) {
						case 500:
							setError('Internal server error');
							break;
						case 405:
							setError('Method not allowed');
							break;
						case 404:
							setError('Resource not found');
							break;
						default:
							setError('Something went wrong adding new todo');
							break;
					}
				}
				const responseJson = await response.json();
				setTodos(responseJson);
			} catch (error) {
				setError('Something went wrong fetching todos');
			}
		}
		fetchTodos();
	}, []);

	async function addTodoHandler(e) {
		e.preventDefault();

		try {
			const response = await fetch('http://localhost:3001/addTodo', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: newTodo,
			});
			const responseJson = await response.json();

			if (!response.ok) {
				switch (response.status) {
					case 500:
						setError('Internal server error');
						return;
					case 405:
						setError('Method not allowed');
						return;
					case 404:
						setError('Resource not found');
						return;
					default:
						setError('Something went wrong adding new todo');
						return;
				}
			}
			setNewTodo('');
			setTodos(responseJson);
		} catch {
			setError('Something went wrong adding new todo');
		}
	}

	async function deleteTodoHandler(todoId) {
		try {
			const response = await fetch('http://localhost:3001/deleteTodo', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: todoId,
			});
			const responseJson = await response.json();

			if (!response.ok) {
				switch (response.status) {
					case 500:
						setError('Internal server error');
						return;
					case 405:
						setError('Method not allowed');
						return;
					case 404:
						setError('Resource not found');
						return;
					default:
						setError('Something went wrong deleting review');
						return;
				}
			}
			setTodos(responseJson);
		} catch {
			setError('Something went wrong deleting review');
		}
	}

	return (
		<>
			<div className="flex h-60 items-center justify-center">
				<h2 className="text-4xl">Todo List</h2>
			</div>

			<div className="grid place-items-center">
				<div>
					{error && <p className="text-red-600">{error}</p>}
					<form
						onSubmit={addTodoHandler}
						className="my-2 flex w-96 items-center justify-between rounded-sm bg-slate-200 p-5 focus:outline-none"
					>
						<input
							type="text"
							placeholder="New Todo"
							className="flex items-center justify-between rounded-sm bg-slate-200 focus:outline-none"
							onChange={(e) => setNewTodo(e.target.value)}
							value={newTodo}
							required
						/>
						<button type="submit" className="text-xl text-gray-400">
							<HiPlus />
						</button>
					</form>
					<ul>
						{todos.map((element, index) => (
							<li
								key={index}
								className="my-2 flex w-96 items-center justify-between break-words rounded-sm bg-slate-200 p-5 text-left"
							>
								<p className="w-1/2">{element.todo}</p>
								<p className=" text-gray-400">
									{getDateTimeFromTimeStamp(element.timestamp)}
								</p>
								<button
									className="text-xl text-red-600"
									onClick={() => {
										deleteTodoHandler(element.id);
									}}
								>
									<HiOutlineTrash />
								</button>
							</li>
						))}
					</ul>
				</div>
			</div>
		</>
	);
}
