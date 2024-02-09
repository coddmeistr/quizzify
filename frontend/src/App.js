import Content from './components/Content/Content';
import Header from './components/Header/Header.js';
import styles from './App.module.css'

function App() {
  return (
    <div className="App">

      <div className={styles.header}>
        <Header />
      </div>

      <div className={styles.content}>
        <Content />
      </div>
      
    </div>
  );
}

export default App;
