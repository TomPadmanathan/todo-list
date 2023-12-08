// React
import { useState, useEffect } from 'react';

// CSS
import './App.css';

// Assets
import { HiOutlineTrash } from 'react-icons/hi';

export default function App() {
	const [todos, setTodos] = useState([]);
	const [error, setError] = useState('');
	const [newTodo, setNewTodo] = useState('');

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

	async function handleSubmit(e) {
		e.preventDefault();
		try {
			const response = await fetch('http://localhost:3001/todos', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({
					newTodo,
				}),
			});
			const responseJson = await response.json();

			if (!response.ok) {
				switch (response.status) {
					case 200:
						setTodos(responseJson);
						break;
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
		} catch {
			setError('Something went wrong adding new todo');
		}
	}

	return (
		<div className="App grid h-screen place-items-center">
			<div className="text-center">
				{error && <p>{error}</p>}
				<h2 className="text-4xl ">Todo List</h2>
				<form
					onSubmit={handleSubmit}
					className="my-2 flex w-96 items-center justify-between rounded-sm bg-slate-200 p-5 focus:outline-none"
				>
					<input
						type="text"
						placeholder="New Todo"
						className="flex items-center justify-between rounded-sm bg-slate-200 focus:outline-none"
						onChange={(e) => setNewTodo(e.target.value)}
						required
					/>
					<button type="submit">Add todo</button>
				</form>
				<ul>
					{todos.map((element, index) => (
						<li
							key={index}
							className="my-2 flex w-96 items-center justify-between rounded-sm bg-slate-200 p-5"
						>
							<p>{element.todo}</p>
							<button className=" text-xl text-red-600" onClick={() => {}}>
								<HiOutlineTrash />
							</button>
						</li>
					))}
				</ul>
			</div>
		</div>
	);
}
