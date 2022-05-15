import React, { useEffect, useState } from 'react';
import axios from 'axios';
import './App.css';

function App() {
  const [file, setFile] = useState<File | null>(null);
  const [progress, setProgress] = useState<number | null>(null);

  const onChange: React.ChangeEventHandler<HTMLInputElement> = (e) => {
    console.log(e.target.files);
    if (e.target.files!.length > 0) {
      setFile(e.target.files?.item(0) || null);
    }
  };

  const onSubmit = async () => {
    if (file) {
      console.log(file);

      const CHANL_SIZE = 1024 * 1024;
      let cur = 0;
      let end = false;

      const ws = new WebSocket(`ws://localhost:1323/api/ws`);
      ws.onmessage = function(evt) {
        if (cur < file.size) {
          let end = cur + CHANL_SIZE;
          if (end > file.size) end = file.size;

          const blob = file.slice(cur, end);
          ws.send(blob);
          cur = end;
          setProgress(Math.floor(cur / file.size * 100));
        } else if (!end) {
          ws.send('');
          end = true;
        } else {
          ws.close();
        }
      }
    }
  };

  return (
    <div className="App">
      <header className="App-header">
        <p>
          <input type="file"  onChange={onChange} />
        </p>
        <button onClick={onSubmit}>UPLOAD</button>
        { progress !== null && (
          <div>loading... {progress}%</div>
        )}
      </header>
    </div>
  );
}

export default App;
