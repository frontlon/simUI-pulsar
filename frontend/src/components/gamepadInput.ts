export type GamepadDirection = 'ArrowUp' | 'ArrowDown' | 'ArrowLeft' | 'ArrowRight'

export interface GamepadHandlers {
  onNavigate: (direction: GamepadDirection) => void
  onConfirm?: () => void
  onRunGame?: () => void
  onPlatformNav?: (direction: 'ArrowUp' | 'ArrowDown') => void
  onMenuNav?: (direction: 'ArrowUp' | 'ArrowDown') => void
  onCancel?: () => void
  onOpenPlatform?: () => void
}

export interface GamepadControllerOptions {
  handlers: GamepadHandlers
  onDisconnect?: () => void
  deadzone?: number
  initialDelay?: number
  repeatInterval?: number
}

const DEFAULT_DEADZONE = 0.35
const DEFAULT_INITIAL_DELAY = 400
const DEFAULT_REPEAT_INTERVAL = 75

class RepeatKey {
  private held = false
  private holdStart = 0
  private lastFire = 0

  constructor(
    private initialDelay: number,
    private repeatInterval: number,
  ) {}

  reset() {
    this.held = false
  }

  update(isDown: boolean, now: number, onFire: () => void) {
    if (!isDown) {
      this.held = false
      return
    }
    if (!this.held) {
      this.held = true
      this.holdStart = now
      this.lastFire = now
      onFire()
      return
    }
    const elapsed = now - this.holdStart
    const interval = elapsed < this.initialDelay ? this.initialDelay : this.repeatInterval
    if (now - this.lastFire >= interval) {
      this.lastFire = now
      onFire()
    }
  }
}

function stickToDirection(x: number, y: number, deadzone: number): GamepadDirection | null {
  const ax = Math.abs(x)
  const ay = Math.abs(y)
  if (ax < deadzone && ay < deadzone) {
    return null
  }
  if (ax > ay) {
    return x < 0 ? 'ArrowLeft' : 'ArrowRight'
  }
  return y < 0 ? 'ArrowUp' : 'ArrowDown'
}

function isPressed(buttons: ReadonlyArray<GamepadButton>, index: number): boolean {
  return !!buttons[index]?.pressed
}

function wasJustPressed(buttons: ReadonlyArray<GamepadButton>, prev: boolean[], index: number): boolean {
  return isPressed(buttons, index) && !prev[index]
}

export function createGamepadController(options: GamepadControllerOptions) {
  const handlers = options.handlers
  const deadzone = options.deadzone ?? DEFAULT_DEADZONE
  const initialDelay = options.initialDelay ?? DEFAULT_INITIAL_DELAY
  const repeatInterval = options.repeatInterval ?? DEFAULT_REPEAT_INTERVAL

  let rafId = 0
  let padIndex = -1
  let prevPressed: boolean[] = []
  let wasConnected = false

  const repeatKeys = {
    up: new RepeatKey(initialDelay, repeatInterval),
    down: new RepeatKey(initialDelay, repeatInterval),
    left: new RepeatKey(initialDelay, repeatInterval),
    right: new RepeatKey(initialDelay, repeatInterval),
    platformUp: new RepeatKey(initialDelay, repeatInterval),
    platformDown: new RepeatKey(initialDelay, repeatInterval),
    menuUp: new RepeatKey(initialDelay, repeatInterval),
    menuDown: new RepeatKey(initialDelay, repeatInterval),
  }

  function resetRepeatKeys() {
    Object.values(repeatKeys).forEach((key) => key.reset())
  }

  function stopLoop() {
    if (rafId) {
      cancelAnimationFrame(rafId)
      rafId = 0
    }
    resetRepeatKeys()
    prevPressed = []
  }

  function tick() {
    rafId = requestAnimationFrame(tick)

    if (!document.hasFocus()) {
      return
    }

    const gamepad = padIndex >= 0 ? navigator.getGamepads()[padIndex] : null

    if (!gamepad) {
      if (wasConnected) {
        wasConnected = false
        options.onDisconnect?.()
      }
      return
    }

    wasConnected = true
    const now = performance.now()
    const buttons = gamepad.buttons
    const axes = gamepad.axes

    const leftDir = stickToDirection(axes[0] ?? 0, axes[1] ?? 0, deadzone)
    const rightDir = stickToDirection(axes[2] ?? 0, axes[3] ?? 0, deadzone)

    const up = isPressed(buttons, 12) || leftDir === 'ArrowUp' || rightDir === 'ArrowUp'
    const down = isPressed(buttons, 13) || leftDir === 'ArrowDown' || rightDir === 'ArrowDown'
    const left = isPressed(buttons, 14) || leftDir === 'ArrowLeft' || rightDir === 'ArrowLeft'
    const right = isPressed(buttons, 15) || leftDir === 'ArrowRight' || rightDir === 'ArrowRight'

    repeatKeys.up.update(up, now, () => handlers.onNavigate('ArrowUp'))
    repeatKeys.down.update(down, now, () => handlers.onNavigate('ArrowDown'))
    repeatKeys.left.update(left, now, () => handlers.onNavigate('ArrowLeft'))
    repeatKeys.right.update(right, now, () => handlers.onNavigate('ArrowRight'))

    if (handlers.onPlatformNav) {
      repeatKeys.platformUp.update(isPressed(buttons, 4), now, () => handlers.onPlatformNav!('ArrowUp'))
      repeatKeys.platformDown.update(isPressed(buttons, 6), now, () => handlers.onPlatformNav!('ArrowDown'))
    }

    if (handlers.onMenuNav) {
      repeatKeys.menuUp.update(isPressed(buttons, 5), now, () => handlers.onMenuNav!('ArrowUp'))
      repeatKeys.menuDown.update(isPressed(buttons, 7), now, () => handlers.onMenuNav!('ArrowDown'))
    }

    if (wasJustPressed(buttons, prevPressed, 0)) {
      handlers.onConfirm?.()
    }
    if (wasJustPressed(buttons, prevPressed, 1)) {
      handlers.onCancel?.()
    }
    if (wasJustPressed(buttons, prevPressed, 2)) {
      handlers.onOpenPlatform?.()
    }
    if (wasJustPressed(buttons, prevPressed, 3)) {
      handlers.onRunGame?.()
    }
    if (wasJustPressed(buttons, prevPressed, 9)) {
      if (handlers.onRunGame) {
        handlers.onRunGame()
      } else {
        handlers.onConfirm?.()
      }
    }

    prevPressed = buttons.map((b) => b.pressed)
  }

  function startLoop(index: number) {
    padIndex = index
    if (!rafId) {
      rafId = requestAnimationFrame(tick)
    }
  }

  function onConnect(e: GamepadEvent) {
    startLoop(e.gamepad.index)
  }

  function onDisconnect(e: GamepadEvent) {
    if (e.gamepad.index === padIndex) {
      stopLoop()
      padIndex = -1
      if (wasConnected) {
        wasConnected = false
        options.onDisconnect?.()
      }
    }
  }

  function init() {
    window.addEventListener('gamepadconnected', onConnect)
    window.addEventListener('gamepaddisconnected', onDisconnect)

    const pads = navigator.getGamepads()
    for (const pad of pads) {
      if (pad) {
        startLoop(pad.index)
        break
      }
    }
  }

  function destroy() {
    window.removeEventListener('gamepadconnected', onConnect)
    window.removeEventListener('gamepaddisconnected', onDisconnect)
    stopLoop()
    padIndex = -1
    wasConnected = false
  }

  return { init, destroy }
}
