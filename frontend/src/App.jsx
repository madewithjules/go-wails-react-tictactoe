import {useState, useEffect} from 'react';
import './App.css';
import {NewGame, MakeMove, ResetGame} from '../wailsjs/go/main/App';

function App() {
  const [board, setBoard] = useState(Array(9).fill(''));
  const [currentPlayer, setCurrentPlayer] = useState('X');
  const [statusMessage, setStatusMessage] = useState("Player X's turn");
  const [theme, setTheme] = useState('light');
  const [isGameOver, setIsGameOver] = useState(false);

  useEffect(() => {
    NewGame()
      .then(initialState => {
        setBoard(initialState.board); // Ensure field names match Go struct (Board vs board)
        setCurrentPlayer(initialState.currentPlayer);
        setStatusMessage(`Player ${initialState.currentPlayer}'s turn`);
        setIsGameOver(initialState.gameOver);
      })
      .catch(err => {
        console.error('Error fetching initial game state:', err);
        setStatusMessage('Error loading game. Please try refreshing.');
      });
  }, []);

  const handleCellClick = index => {
    if (isGameOver || board[index] !== '') {
      return;
    }

    MakeMove(index)
      .then(updatedState => {
        setBoard(updatedState.board);
        setCurrentPlayer(updatedState.currentPlayer);
        setIsGameOver(updatedState.gameOver);

        if (updatedState.winner) {
          if (updatedState.winner === 'draw') {
            setStatusMessage("It's a draw!");
          } else {
            setStatusMessage(`Player ${updatedState.winner} wins!`);
          }
          // Auto-reset after a delay
          setTimeout(() => {
            handleResetClick();
          }, 2000);
        } else {
          setStatusMessage(`Player ${updatedState.currentPlayer}'s turn`);
        }
      })
      .catch(errorMessage => {
        setStatusMessage(String(errorMessage)); // Ensure errorMessage is a string
      });
  };

  const renderBoard = () => {
    return (
      <div className="board">
        {board.map((cell, index) => (
          <button
            key={index}
            className="cell"
            onClick={() => handleCellClick(index)}
            disabled={isGameOver || cell !== ''}
          >
            {cell}
          </button>
        ))}
      </div>
    );
  };

  const handleResetClick = () => {
    ResetGame()
      .then(newState => {
        setBoard(newState.board);
        setCurrentPlayer(newState.currentPlayer);
        setStatusMessage(`Player ${newState.currentPlayer}'s turn`);
        setIsGameOver(newState.gameOver);
      })
      .catch(err => {
        console.error('Error resetting game:', err);
        setStatusMessage('Error resetting game. Please try again.');
      });
  };

  const toggleTheme = () => {
    setTheme(currentTheme => (currentTheme === 'light' ? 'dark' : 'light'));
  };

  return (
    <div id="App" className={theme === 'dark' ? 'dark-theme' : ''}>
      <div className="game-container">
        <h1>Tic-Tac-Toe</h1>
        {renderBoard()}
        <div className="status-message">{statusMessage}</div>
        <div className="controls">
          <button className="btn" onClick={handleResetClick}>
            Reset Game
          </button>
          <button className="btn theme-toggle" onClick={toggleTheme}>
            Switch to {theme === 'light' ? 'Dark' : 'Light'} Theme
          </button>
        </div>
      </div>
    </div>
  );
}

export default App;
