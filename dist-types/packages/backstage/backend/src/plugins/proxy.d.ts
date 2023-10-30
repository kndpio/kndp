import { Router } from 'express';
import { PluginEnvironment } from '../types';
export default function createPlugin(env: PluginEnvironment): Promise<Router>;
