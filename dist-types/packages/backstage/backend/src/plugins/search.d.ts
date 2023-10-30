import { PluginEnvironment } from '../types';
import { Router } from 'express';
export default function createPlugin(env: PluginEnvironment): Promise<Router>;
