interface Props {
  SYSTEM_NAME: string;
  USER_NAME: string;
  USER_LAST_NAME: string;
  VERIFICATION_URL: string;
}

export function getVerifyAccountHTML({
  SYSTEM_NAME,
  USER_NAME,
  USER_LAST_NAME,
  VERIFICATION_URL,
}: Props) {
  return `<!doctype html>
  <html lang="es">
    <head>
      <meta charset="UTF-8" />
      <meta name="viewport" content="width=device-width, initial-scale=1.0" />
      <title>${SYSTEM_NAME} - Confirma tu cuenta</title>
      <style>
        @import url('https://fonts.googleapis.com/css2?family=Poppins:ital,wght@0,100;0,200;0,300;0,400;0,500;0,600;0,700;0,800;0,900;1,100;1,200;1,300;1,400;1,500;1,600;1,700;1,800;1,900&display=swap');
      </style>
      <style>
        body {
          background-color: #f6f6f6;
          font-family: 'Poppins', sans-serif;
          font-weight: 400;
          font-style: normal;
          font-size: 16px;
          color: #334155;
          margin: 0;
          padding: 0;
        }
        h1 {
          text-align: center;
          color: #059669;
          margin-top: 0;
        }
        p {
          margin: 0;
        }
        a {
          all: unset;
        }
        .wrapper {
          height: 100vh;
          display: grid;
          place-items: center;
        }
        .container {
          max-width: 600px;
          margin: 0 auto;
          padding: 20px;
          background-color: #fff;
          border-radius: 10px;
  
          display: flex;
          flex-direction: column;
          gap: 20px;
          align-items: center;
        }
        .content_mail {
          width: 100%;
        }
        .button {
          all: unset;
          cursor: pointer;
          display: inline-block;
          padding: 10px 20px;
  
          background-color: #10b981;
          color: #fff;
          text-decoration: none;
          border-radius: 5px;
          transition: 0.2s;
  
          &:hover {
            background-color: #059669;
          }
        }
  
        .mini_content {
          font-size: 13px;
          color: #a1a1aa;
        }
        .footer {
          display: flex;
          flex-direction: column;
          align-items: center;
        }
        .alert_message {
          font-size: 13px;
          color: #a1a1aa;
        }
      </style>
    </head>
    <body>
      <div class="wrapper">
        <div class="container">
          <h1>Te damos la bienvenida a ${SYSTEM_NAME} 👋</h1>
  
          <div class="content_mail">
            <p>Hola ${USER_NAME} ${USER_LAST_NAME}</p>
            <p>Confirma tu dirección de correo para completar tu registro:</p>
          </div>
  
          <a href="${VERIFICATION_URL}">
            <button class="button">
              Confirmar mi dirección de correo electrónico
            </button>
          </a>
  
          <p class="mini_content">Mensaje de contenido</p>
  
          <div class="footer">
            <p class="alert_message">
              Si no creaste esta cuenta, por favor ignora este correo electrónico.
            </p>
            <p class="alert_message">
              Este es un servicio de notificación por correo.
            </p>
          </div>
        </div>
      </div>
    </body>
  </html>
  
`;
}
