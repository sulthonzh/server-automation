<?php

namespace App\Http\Controllers;

use Illuminate\Http\Request;
use Symfony\Component\Process\Process;
use Symfony\Component\Process\Exception\ProcessFailedException;

class PhpFpmController extends Controller
{
    public function setPhpFpm(Request $request)
    {
        $version = $request->input('version');
        $dir = $request->input('dir');
        $goAgentPath = env('GO_AGENT_PATH', './go-agent');

        $process = new Process([$goAgentPath, "phpfpm", "--version={$version}", "--dir={$dir}"]);
        $process->run();

        if (!$process->isSuccessful()) {
            return response()->json(['success' => false, 'message' => 'Failed to PHP-FPM.'], 500);
        }

        return response()->json(['success' => true, 'message' => 'PHP-FPM set successfully.']);
    }
}
