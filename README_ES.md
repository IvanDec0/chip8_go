# Emulador Chip-8

Un emulador Chip-8 moderno y configurable escrito en Go con SDL2, que presenta un manejo robusto de errores, seguridad de memoria y configuraciones personalizables de visualización/rendimiento.

## Tabla de Contenidos

- [Visión General](#visión-general)
- [Características](#características)
- [Prerrequisitos](#prerrequisitos)
- [Instalación](#instalación)
- [Uso](#uso)
- [Opciones de Configuración](#opciones-de-configuración)
- [Controles](#controles)
- [Estructura del Proyecto](#estructura-del-proyecto)
- [Detalles Técnicos](#detalles-técnicos)
- [Compatibilidad con ROMs](#compatibilidad-con-roms)
- [Solución de Problemas](#solución-de-problemas)

## Visión General

El Chip-8 es un lenguaje de programación interpretado desarrollado en la década de 1970 para su uso en microcomputadoras de 8 bits. Este emulador recrea fielmente el sistema Chip-8, permitiéndote ejecutar juegos y programas clásicos diseñados originalmente para la plataforma.

### ¿Qué es Chip-8?

- **Plataforma**: Originalmente diseñado para la computadora COSMAC VIP
- **Pantalla**: Pantalla monocromática de 64×32 píxeles
- **Memoria**: 4KB de RAM
- **Registros**: 16 registros de propósito general de 8 bits (V0-VF)
- **Entrada**: Teclado hexadecimal de 16 teclas
- **Sonido**: Pitido de un solo tono (440Hz)

## Características

### Emulación Básica

- ✅ Implementación completa del conjunto de instrucciones Chip-8
- ✅ Temporización precisa y renderizado de pantalla
- ✅ Soporte para temporizador de sonido (indicación de pitido)
- ✅ Manejo de entrada de 16 teclas
- ✅ Operación segura de memoria con verificación de límites

### Mejoras Modernas

- ✅ **Escalado de resolución configurable** (1x a 20x)
- ✅ **Velocidad de CPU ajustable** (100-2000 ciclos/segundo)
- ✅ **Títulos de ventana personalizables**
- ✅ **Soporte de sonido** con audio SDL2 (pitido de 440Hz)
- ✅ **Arquitectura limpia** con separación de responsabilidades
- ✅ **Diseño basado en interfaces** para extensibilidad
- ✅ **Manejo robusto de errores** y validación de entrada
- ✅ **Protección de memoria** contra desbordamientos de búfer
- ✅ **Protección de pila**: Prevención de condiciones de desbordamiento/subdesbordamiento
- ✅ **Soporte multiplataforma** (Linux, Windows, macOS)

### Control de Calidad

- ✅ Verificación exhaustiva de límites para todas las operaciones de memoria
- ✅ Renderizado seguro de sprites con validación de coordenadas
- ✅ Mecanismos protegidos de llamada/retorno de subrutinas
- ✅ Validación de parámetros de entrada con mensajes de error claros

## Prerrequisitos

- **Go**: Versión 1.19 o posterior
- **SDL2**: Bibliotecas de desarrollo
  - **Ubuntu/Debian**: `sudo apt-get install libsdl2-dev`
  - **Fedora/CentOS**: `sudo dnf install SDL2-devel`
  - **macOS**: `brew install sdl2`
  - **Windows**: Descargar desde [SDL2 releases](https://github.com/libsdl-org/SDL/releases)

## Instalación

1. **Clonar el repositorio:**

   ```bash
   git clone https://github.com/IvanDec0/chip8_go
   cd chip8
   ```

2. **Instalar dependencias:**

   ```bash
   go mod download
   ```

3. **Compilar el emulador:**
   ```bash
   go build -o chip8-emulator
   ```

## Uso

### Uso Básico

```bash
# Ejecutar con configuración predeterminada
./chip8-emulator <archivo-ROM>

# Ejemplo
./chip8-emulator roms/games/Pong.ch8
```

### Uso Avanzado

```bash
# Configuración personalizada
./chip8-emulator [opciones] <archivo-ROM>

# Ejemplos
./chip8-emulator -scale 15 -speed 500 -title "Juego Pong" roms/games/Pong.ch8
./chip8-emulator -scale 5 -speed 200 roms/demos/Maze\ \[David\ Winter,\ 199x\].ch8
```

### Ayuda

```bash
./chip8-emulator --help
```

## Opciones de Configuración

| Bandera  | Descripción                     | Rango            | Predeterminado |
| -------- | ------------------------------- | ---------------- | -------------- |
| `-scale` | Factor de escala de ventana     | 1-20             | 10             |
| `-speed` | Ciclos de CPU por segundo       | 100-2000         | 700            |
| `-title` | Título personalizado de ventana | Cualquier cadena | "Chip8"        |

### Ejemplos de Escala

- **Escala 1**: 64×32 píxeles (tamaño original)
- **Escala 10**: 640×320 píxeles (predeterminado)
- **Escala 20**: 1280×640 píxeles (máximo)

### Ejemplos de Velocidad

- **100 ciclos/seg**: Muy lento, bueno para depuración
- **700 ciclos/seg**: Velocidad predeterminada, funciona bien para la mayoría de juegos
- **2000 ciclos/seg**: Ejecución rápida para demos

## Controles

El emulador mapea tu teclado al teclado hexadecimal de 16 teclas del Chip-8:

```
Teclado Chip-8    Mapeo de Teclado
┌─┬─┬─┬─┐        ┌─┬─┬─┬─┐
│1│2│3│C│        │1│2│3│4│
├─┼─┼─┼─┤        ├─┼─┼─┼─┤
│4│5│6│D│        │Q│W│E│R│
├─┼─┼─┼─┤   =>   ├─┼─┼─┼─┤
│7│8│9│E│        │A│S│D│F│
├─┼─┼─┼─┤        ├─┼─┼─┼─┤
│A│0│B│F│        │Z│X│C│V│
└─┴─┴─┴─┘        └─┴─┴─┴─┘
```

### Teclas Especiales

- **ESC**: Salir del emulador
- **Cerrar Ventana (X)**: Salir del emulador

## Características Avanzadas

### Sistema de Sonido

El emulador incluye una implementación completa de audio:

- **Frecuencia**: 440Hz (nota A4)
- **Reproducción automática**: El sonido se reproduce cuando los programas CHIP-8 usan el temporizador de sonido

### Ejemplos de Configuración

Perfecto para diferentes casos de uso:

- **Desarrollo**: `./chip8-emulator -scale 5 -speed 200 rom.ch8` (pequeño, lento)
- **Juegos**: `./chip8-emulator -scale 15 -speed 700 rom.ch8` (grande, velocidad normal)
- **Demos**: `./chip8-emulator -scale 20 -speed 1000 rom.ch8` (pantalla completa, rápido)

## Estructura del Proyecto

```
chip8/
├── main.go                    # Punto de entrada simple y configuración
├── internal/                  # Paquetes de implementación privados
│   ├── core/                  # Lógica pura de CPU CHIP-8 (sin dependencias)
│   │   ├── cpu.go             # Estado de CPU y funcionalidad principal
│   │   └── opcodes.go         # Implementación del conjunto de instrucciones
│   ├── audio/                 # Interface y implementación del sistema de audio
│   │   ├── audio.go           # Definición de interface de audio
│   │   └── sdl_audio.go       # Implementación de audio SDL2
│   ├── renderer/              # Interface y implementación de renderizado
│   │   ├── renderer.go        # Definición de interface de renderizado
│   │   └── sdl_renderer.go    # Implementación de renderizado SDL2
│   └── input/                 # Manejo de entrada
│       └── handler.go         # Procesamiento de entrada de teclado
├── pkg/                       # Paquetes públicos
│   └── emulator/              # Coordinación de emulador de alto nivel
│       └── emulator.go        # Orquestación principal del emulador
├── roms/                      # Colección de ROMs
│   ├── demos/                 # Programas de demostración
│   ├── games/                 # Juegos clásicos
│   └── programs/              # Programas utilitarios
├── go.mod                     # Dependencias del módulo Go
├── go.sum                     # Sumas de verificación de dependencias
└──  README.md                 # Esta documentación
```

### Componentes Clave

- **`main.go`**: Punto de entrada limpio con análisis de configuración
- **`internal/core/`**: Lógica pura de CPU CHIP-8, independiente de plataforma
- **`internal/audio/`**: Sistema de audio basado en interfaces con implementación SDL2
- **`internal/renderer/`**: Sistema de renderizado basado en interfaces con implementación SDL2
- **`internal/input/`**: Manejo de entrada dirigido por eventos
- **`pkg/emulator/`**: API pública que coordina todos los componentes
- **`roms/`**: Colección de programas y juegos CHIP-8

## Detalles Técnicos

### Disposición de Memoria

```
0x000-0x1FF: Intérprete Chip-8 (datos de fuente almacenados en 0x50-0x9F)
0x200-0xFFF: ROM del programa y RAM (3584 bytes)
```

### Registros

- **V0-VF**: 16 registros de propósito general de 8 bits
- **I**: Registro de dirección de 16 bits
- **PC**: Contador de programa
- **SP**: Puntero de pila
- **DT**: Temporizador de retardo (decrece a 60Hz)
- **ST**: Temporizador de sonido (decrece a 60Hz, emite pitido cuando > 0)

### Pantalla

- **Resolución**: 64×32 píxeles
- **Colores**: Monocromático (negro/blanco)
- **Renderizado**: Dibujo de sprites basado en XOR
- **Refresco**: 60 FPS

### Conjunto de Instrucciones

El emulador implementa todos los 35 opcodes estándar de Chip-8:

- Operaciones de memoria (cargar, almacenar, copiar)
- Operaciones aritméticas y lógicas
- Flujo de control (salto, llamada, retorno, omisión)
- Gráficos (limpiar pantalla, dibujar sprite)
- Manejo de entrada (detección de pulsación de teclas)
- Operaciones de temporizador

### Características de Seguridad

1. **Protección de Memoria**:

   - Verificación de límites para todo acceso a memoria
   - Prevención de desbordamientos de búfer
   - Carga segura de ROM con validación de tamaño

2. **Protección de Pila**:

   - Prevención de desbordamiento de pila (máx. 16 niveles)
   - Protección contra subdesbordamiento de pila
   - Manejo seguro de subrutinas

3. **Seguridad de Pantalla**:
   - Envolvimiento de coordenadas y verificación de límites
   - Renderizado protegido de sprites
   - Manipulación segura de píxeles

## Compatibilidad con ROMs

El emulador es compatible con ROMs estándar de Chip-8 y ha sido probado con:

### Juegos Incluidos

- **Pong**: Juego clásico de paletas
- **Breakout**: Juego de romper ladrillos
- **Tetris**: Juego de puzzle de bloques
- **Space Invaders**: Juego de disparos arcade
- **Pac-Man**: Juego de navegación por laberinto

### Demos Incluidas

- **Maze**: Generación procedimental de laberintos
- **Particle Demo**: Demostración de efectos visuales
- **Sierpinski**: Generación de patrones fractales
- **Stars**: Campo de estrellas animado

### Formatos de Archivo

- **Extensión**: `.ch8` (estándar)
- **Tamaño**: Hasta 3584 bytes
- **Formato**: Datos binarios sin procesar

## Solución de Problemas

### Problemas Comunes

1. **"Error al inicializar SDL"**

   - Asegúrate de que las bibliotecas de desarrollo SDL2 estén instaladas
   - Verifica que tus controladores de gráficos estén actualizados

2. **"Error al cargar ROM"**

   - Verifica que el archivo ROM exista y sea legible
   - Revisa los permisos del archivo
   - Asegúrate de que la ROM sea un archivo Chip-8 válido

3. **Ventana demasiado pequeña/grande**

   - Ajusta el parámetro `-scale` (1-20)
   - Prueba `-scale 10` para un buen tamaño predeterminado

4. **Juego corriendo demasiado rápido/lento**
   - Ajusta el parámetro `-speed` (100-2000)
   - La mayoría de los juegos funcionan bien con 500-1000 ciclos/segundo

### Consejos de Rendimiento

- Usa valores de escala más bajos para un mejor rendimiento en hardware antiguo
- Reduce la velocidad de CPU para ROMs complejas que corren demasiado rápido
- Cierra otras aplicaciones si experimentas tartamudeo

---

## Ejemplos de Inicio Rápido

```bash
# Descargar y ejecutar un juego clásico
./chip8-emulator roms/games/Pong.ch8

# Ejecutar una demo con configuración personalizada
./chip8-emulator -scale 15 -speed 800 -title "Demo Chip8" roms/demos/Maze\ \[David\ Winter,\ 199x\].ch8

# Modo depuración (velocidad lenta, ventana pequeña)
./chip8-emulator -scale 5 -speed 200 tu-rom.ch8
```

¡Disfruta explorando el mundo de la programación y los juegos Chip-8! 🎮
