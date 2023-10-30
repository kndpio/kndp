import { Router } from 'express';
import type { PluginEnvironment } from '../types';
export default function createPlugin(env: PluginEnvironment): Promise<Router>;
