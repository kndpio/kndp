import express, { Request, Response } from 'express';
import bodyParser from 'body-parser';
import fs from 'fs';
import cors from 'cors';
import path from 'path'; 

const app = express();
const port = 3000;

app.use(bodyParser.json());

app.use(cors()); 

app.post('/profile', (req: Request, res: Response) => {
  const { logo, main_color, background_color } = req.body;

  const generatedCSS = `
  .sidebar__logo__character {
    visibility: hidden; 
  }
  .sidebar__logo::before {
    content:"";
    position: absolute;
    left: 20px;
    top: 30px;
    width: 40px; 
    height: 40px; 
    background: url("${logo}");
    background-size: 40px 40px;
  }
    .sidebar {  
      background-color: ${background_color};
    }
    .sidebar__nav-item {
      color: ${main_color};
    }
  `;

  process.env.PATHS = '/shared/app/css';
  
  const outputDir = process.env.PATHS ; 
  const fileName = 'generated.css';
  const filePath = path.join(outputDir, fileName); 

  fs.writeFile(filePath, generatedCSS, (err) => {
    if (err) {
      console.error(err);
      res.status(500).json({ error: 'Error generating and saving CSS' });
    } else {
      console.log('CSS generated and saved successfully!');
      res.json({ success: true });
    }
  });
});

app.listen(port, () => {
  console.log(`Backend API listening at http://localhost:${port}`);
});
