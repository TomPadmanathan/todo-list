// React
import { useState, useEffect, useRef } from 'react';

// CSS
import './App.css';

// Assets
import {
	HiOutlinePencil,
	HiOutlineTrash,
	HiPlus,
	HiCheck,
} from 'react-icons/hi';

export default function App() {
	const [todos, setTodos] = useState([]);
	const [error, setError] = useState('');
	const [newTodo, setNewTodo] = useState('');
	const [isTodoEditing, setIsTodoEditing] = useState(false);

	function checkResponseError(response) {
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
	}

	useEffect(() => {
		async function fetchTodos() {
			try {
				const response = await fetch('http://localhost:3001/todos', {
					method: 'GET',
					headers: { 'Content-Type': 'application/json' },
				});
				checkResponseError(response);
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

			checkResponseError(response);

			setNewTodo('');
			setTodos(responseJson);
		} catch {
			setError('Something went wrong adding new todo');
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
						className="my-2 flex w-[500px] items-center justify-between rounded-sm bg-slate-200 p-5 focus:outline-none"
					>
						<input
							type="text"
							placeholder="New Todo"
							className="mr-2 w-full items-center justify-between rounded-sm bg-slate-200 focus:outline-none"
							onChange={(e) => setNewTodo(e.target.value)}
							value={newTodo}
							required
						/>
						<button
							type="submit"
							className={`text-xl text-gray-400 transition-all hover:text-green-400`}
						>
							<HiPlus />
						</button>
					</form>
					<ul>
						{todos.map((element) => (
							<TodoListItem
								key={element.id}
								element={element}
								setError={setError}
								setTodos={setTodos}
								isTodoEditingState={[isTodoEditing, setIsTodoEditing]}
							/>
						))}
					</ul>
				</div>
			</div>
		</>
	);
}
function TodoListItem({ element, setError, setTodos, isTodoEditingState }) {
	const [isTodoEditing, setIsTodoEditing] = isTodoEditingState;
	const [editable, setEditable] = useState(false);
	const [todoValue, setTodoValue] = useState(element.todo);
	const todoInput = useRef(null);

	useEffect(() => {
		if (isTodoEditing) todoInput.current.focus();
	}, [isTodoEditing]);

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

	async function deleteTodoHandler(todoId) {
		try {
			const response = await fetch('http://localhost:3001/deleteTodo', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: todoId,
			});
			const responseJson = await response.json();

			checkResponseError(response);

			setTodos(responseJson);
		} catch {
			setError('Something went wrong deleting todo');
		}
	}

	async function editTodoHandler(todoId, newValue) {
		try {
			const response = await fetch('http://localhost:3001/editTodo', {
				method: 'POST',
				headers: { 'Content-Type': 'application/json' },
				body: JSON.stringify({ todoId, newValue }),
			});
			const responseJson = await response.json();

			checkResponseError(response);
			setTodos(responseJson);
		} catch {
			setError('Something went wrong editing todo');
		}
	}

	useEffect(() => {
		if (todoInput.current) {
			const element = todoInput.current;
			element.style.height = '25px';
			element.style.height = `${element.scrollHeight}px`;
		}
	}, [todoInput]);

	return (
		<li className="my-2 flex w-[500px] items-center justify-around break-words rounded-sm bg-slate-200 p-5 text-left">
			<textarea
				className="h-[25px] w-1/2 resize-none overflow-y-hidden bg-slate-200 leading-6 focus:outline-none"
				value={todoValue}
				ref={todoInput}
				disabled={!editable}
				onChange={(e) => {
					setTodoValue(e.target.value);

					// Dynamicly change height of element
					e.target.style.height = '25px';
					e.target.style.height = `${e.target.scrollHeight}px`;
				}}
			/>

			<p className="text-gray-400">
				{getDateTimeFromTimeStamp(element.timestamp)}
			</p>

			<button
				className={`text-xl transition-all ${
					!isTodoEditing && 'hover:text-green-400'
				} ${editable ? 'text-green-400' : 'text-gray-400'}`}
				onClick={() => {
					if (editable) {
						editTodoHandler(element.id, todoValue);
						setEditable(false);
						setIsTodoEditing(false);
						return;
					}
					if (isTodoEditing) return;
					setIsTodoEditing(true);
					setEditable(!editable);
				}}
			>
				{!editable ? <HiOutlinePencil /> : <HiCheck />}
			</button>
			<button
				className="text-xl text-red-600"
				onClick={() => deleteTodoHandler(element.id)}
			>
				<HiOutlineTrash />
			</button>
		</li>
	);
}
