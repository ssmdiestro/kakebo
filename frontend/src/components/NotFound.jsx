import errorImg from "../assets/404.png";

export default function NotFound() {
    return (
        <div
            style={{
                height: "100vh",
                display: "flex",
                flexDirection: "column",
                justifyContent: "center",
                alignItems: "center",
                backgroundColor: "#fff",      // Fondo blanco
                textAlign: "center",
            }}
        >
            <img
                src={errorImg}
                alt="Error 404"
                style={{
                    maxWidth: "360px",
                    width: "80%",
                    marginBottom: "24px",
                }}
            />
        </div>
    );
}
