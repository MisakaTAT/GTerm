export interface XTermTheme {
  name: string;
  foreground: string;
  background: string;
  black: string;
  red: string;
  green: string;
  yellow: string;
  blue: string;
  magenta: string;
  cyan: string;
  white: string;
  brightBlack: string;
  brightRed: string;
  brightGreen: string;
  brightYellow: string;
  brightBlue: string;
  brightMagenta: string;
  brightCyan: string;
  brightWhite: string;
}

export const defaultTheme: XTermTheme = {
  name: 'Default',
  foreground: '#ffffff',
  background: '#000000',
  black: '#000000',
  red: '#cd0000',
  green: '#00cd00',
  yellow: '#cdcd00',
  blue: '#0000ee',
  magenta: '#cd00cd',
  cyan: '#00cdcd',
  white: '#e5e5e5',
  brightBlack: '#7f7f7f',
  brightRed: '#ff0000',
  brightGreen: '#00ff00',
  brightYellow: '#ffff00',
  brightBlue: '#5c5cff',
  brightMagenta: '#ff00ff',
  brightCyan: '#00ffff',
  brightWhite: '#ffffff',
};

const themeModules: { [key: string]: () => Promise<any> } = {
  Dracula: () => import('./Dracula'),
  'Solarized Dark': () => import('./Solarized_Dark'),
  Monokai: () => import('./Monokai'),
  'One Dark': () => import('./One_Dark'),
  'One Light': () => import('./One_Light'),
};

export const themeOptions = [
  { label: defaultTheme.name, value: defaultTheme.name },
  ...Object.keys(themeModules).map(name => ({
    label: name,
    value: name,
  })),
].sort((a, b) => a.label.localeCompare(b.label));

export const loadTheme = async (themeName: string): Promise<XTermTheme> => {
  if (themeName === defaultTheme.name) {
    return defaultTheme;
  }

  const module = themeModules[themeName];
  if (!module) {
    console.warn(`Theme ${themeName} not found in themeModules`);
    return defaultTheme;
  }

  try {
    const theme = await module();
    return theme.default;
  } catch {
    return defaultTheme;
  }
};
