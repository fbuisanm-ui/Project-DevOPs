// ============================================================
// CONFIGURACIÓN DEL MÓDULO DE GO
// ============================================================

// MODULE define el nombre o identificador de nuestro módulo.
//
// Un módulo en Go es la unidad que permite organizar,
// desarrollar y gestionar un proyecto.
//
// El nombre "project-devops" identifica nuestro proyecto.
// También se utiliza como referencia cuando otros paquetes
// del proyecto necesitan importar código entre sí.
//
// En proyectos publicados, normalmente se utiliza una ruta
// como "github.com/usuario/nombre-del-proyecto".
module project-devops

// GO establece la versión del lenguaje Go que utiliza el módulo.
//
// En este caso:
// - 1.24.0 indica la versión de Go requerida por el proyecto.
// - Ayuda a Go a determinar qué características del lenguaje
//   y reglas de compatibilidad debe aplicar.
go 1.24.0