#include <WiFi.h>
#include <HTTPClient.h>

// WiFi network name and password:
const char * networkName = "UofM-Guest";
const char * networkPswd = "";

// Internet domain to request from:
const char * hostDomain = "trana.fit";
const int hostPort = 443;

const char* serverNameTrial = "https://trana.fit/newTrial";
const char* serverNameData = "https://trana.fit/newData";

const char* root_ca= \
"-----BEGIN CERTIFICATE-----\n" \
"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\n" \
"-----END CERTIFICATE-----\n";

const int BUTTON_PIN = 0;
const int LED_PIN = 13;

int ledState = 0;
int trialID = -1;
int period = 20;
unsigned long time_now = 0;
int count = 0;
String jsonData = "[";

WiFiClientSecure client;
HTTPClient http;
  
void setup()
{
  // Initilize hardware:
  Serial.begin(115200);
  pinMode(BUTTON_PIN, INPUT_PULLUP);
  pinMode(LED_PIN, OUTPUT);

  // Connect to the WiFi network (see function below loop)
  connectToWiFi(networkName, networkPswd);
  client.setCACert(root_ca);

  digitalWrite(LED_PIN, LOW); // LED off
  Serial.print("Press button 0 to connect to ");
  Serial.println(hostDomain);

}

void loop()
{

  if (digitalRead(BUTTON_PIN) == LOW)
  { // Check if button has been pressed
    while (digitalRead(BUTTON_PIN) == LOW)
      ; // Wait for button to be released
    if (trialID == -1) {
      trialID = httpGETRequest(serverNameTrial).toInt(); // Connect to server
      ledState = 1;
      time_now = millis();
    } else {
      uint16_t value1 = analogRead(32);
      uint16_t value2 = analogRead(33);
      uint16_t value3 = analogRead(34);
      uint16_t value4 = analogRead(39);
      uint16_t value5 = analogRead(36);
      jsonData += "{\"trialid\":" + String(trialID) + ",\"time\":" + String(millis()) + ",\"a1\":" + String(value1) + ",\"a2\":" + String(value2) + ",\"a3\":" + String(value3) + ",\"a4\":" + String(value4) + ",\"a5\":" + String(value5) + "}]";
//      Serial.println(jsonData);
      httpPOSTRequest(serverNameData, jsonData);
      count = 0;
      jsonData = "[";
      trialID = -1;
      ledState = 0;
    }
    digitalWrite(LED_PIN, ledState); // Turn off LED
    Serial.println(trialID);
  } else {
    uint16_t value1 = analogRead(32);
    uint16_t value2 = analogRead(33);
    uint16_t value3 = analogRead(34);
    uint16_t value4 = analogRead(39);
    uint16_t value5 = analogRead(36);
    if ((trialID != -1) && (millis() > time_now + period)) {
      if (count == 49) {
        jsonData += "{\"trialid\":" + String(trialID) + ",\"time\":" + String(millis()) + ",\"a1\":" + String(value1) + ",\"a2\":" + String(value2) + ",\"a3\":" + String(value3) + ",\"a4\":" + String(value4) + ",\"a5\":" + String(value5) + "}]";
        httpPOSTRequest(serverNameData, jsonData);
        count = 0;
        jsonData = "[";
      } else {
        jsonData += "{\"trialid\":" + String(trialID) + ",\"time\":" + String(millis()) + ",\"a1\":" + String(value1) + ",\"a2\":" + String(value2) + ",\"a3\":" + String(value3) + ",\"a4\":" + String(value4) + ",\"a5\":" + String(value5) + "},";
//        Serial.println(jsonData);
        count++;
      }
      time_now = millis();
    }
  }

}

void connectToWiFi(const char * ssid, const char * pwd)
{
  printLine();
  Serial.println("Connecting to WiFi network: " + String(ssid));

  WiFi.begin(ssid, pwd);

  while (WiFi.status() != WL_CONNECTED)
  {
    // Blink LED while we're connecting:
    digitalWrite(LED_PIN, ledState);
    ledState = (ledState + 1) % 2; // Flip ledState
    delay(500);
    Serial.print(".");
  }

  Serial.println();
  Serial.println("WiFi connected!");
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());

}

String httpGETRequest(const char* serverName) {

  // Your Domain name with URL path or IP address with path
  http.begin(client, serverName, hostPort);

  // Send HTTP POST request
  int httpResponseCode = http.GET();

  String payload = "--";

  if (httpResponseCode > 0) {
    Serial.print("HTTP Response code: ");
    Serial.println(httpResponseCode);
    payload = http.getString();
  }
  else {
    Serial.print("Error code: ");
    Serial.println(httpResponseCode);
  }
  // Free resources
  http.end();

  return payload;
}

int httpPOSTRequest(const char* serverName, String Data) {

  // Your Domain name with URL path or IP address with path
  http.begin(client, serverName, hostPort);

  // If you need an HTTP request with a content type: application/json, use the following:
  http.addHeader("Content-Type", "application/json");
  int httpResponseCode = http.POST(Data);

  // Free resources
  http.end();

  return httpResponseCode;
}

void printLine()
{
  Serial.println();
  for (int i = 0; i < 30; i++)
    Serial.print("-");
  Serial.println();
}
